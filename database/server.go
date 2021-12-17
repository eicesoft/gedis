package database

import (
	"fmt"
	"gedis/aof"
	"gedis/config"
	"gedis/pkg/logger"
	"gedis/pkg/utils"
	"gedis/reply"
	"gedis/types/cmd"
	"gedis/types/pubsub"
	"gedis/types/redis"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	DbSize = 16
)

type MultiDB struct {
	dbSet      []*DB
	hub        *pubsub.Hub //public/subscribe handle
	aofHandler *aof.Handler
}

func MakeBasicMultiDB() *MultiDB {
	mdb := &MultiDB{}
	mdb.dbSet = make([]*DB, DbSize)
	for i := range mdb.dbSet {
		mdb.dbSet[i] = makeBasicDB()
	}
	return mdb
}

// NewStandaloneServer 构造一个单实例服务器
func NewStandaloneServer() *MultiDB {
	mdb := &MultiDB{}
	mdb.dbSet = make([]*DB, DbSize)
	for i := range mdb.dbSet {
		singleDB := makeDB()
		singleDB.index = int8(i)
		mdb.dbSet[i] = singleDB
	}
	mdb.hub = pubsub.MakeHub()
	if config.Get().Server.EnableAof {
		aofHandler, err := aof.NewAOFHandler(mdb, func() redis.EmbedDB {
			return MakeBasicMultiDB()
		})
		if err != nil {
			panic(err)
		}
		mdb.aofHandler = aofHandler
	}

	for _, db := range mdb.dbSet {
		singleDB := db
		singleDB.addAof = func(line CmdLine) { //注册Aof写入函数
			mdb.aofHandler.AddAof(int(singleDB.index), line)
		}
	}
	return mdb
}

func (mdb *MultiDB) Exec(c redis.Connection, cmdLine [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &reply.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(cmdLine[0]))

	if cmdName == cmd.Auth { // authenticate 密码
		return Auth(c, cmdLine[1:])
	}

	if !isAuthenticated(c) { //authenticate 是否通过
		return reply.MakeErrReply("NOAUTH Authentication required")
	}

	if cmdName == cmd.Subscribe { // scribe commands handle
		if len(cmdLine) < 2 {
			return reply.MakeArgNumErrReply("subscribe")
		}
		return pubsub.Subscribe(mdb.hub, c, cmdLine[1:])
	} else if cmdName == cmd.Publish {
		return pubsub.Publish(mdb.hub, cmdLine[1:])
	} else if cmdName == cmd.UnSubscribe {
		return pubsub.UnSubscribe(mdb.hub, c, cmdLine[1:])
	} else if cmdName == cmd.BgRewriteAof {
		// aof.go imports router.go, router.go cannot import BGRewriteAOF from aof.go
		return BGRewriteAOF(mdb, cmdLine[1:])
	} else if cmdName == cmd.RewriteAof {
		return RewriteAOF(mdb, cmdLine[1:])
	} else if cmdName == cmd.FlushAll {
		return mdb.flushAll()
	} else if cmdName == cmd.Select {
		if c != nil && c.InMultiState() { //multi指令不能切换数据库
			return reply.MakeErrReply("cannot select database within multi")
		}
		if len(cmdLine) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(c, mdb, cmdLine[1:])
	}

	dbIndex := c.GetDBIndex() // normal commands

	if dbIndex >= len(mdb.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	selectedDB := mdb.dbSet[dbIndex]
	return selectedDB.Exec(c, cmdLine)
}

// AfterClientClose does some clean after client close connection
func (mdb *MultiDB) AfterClientClose(c redis.Connection) {
	pubsub.UnsubscribeAll(mdb.hub, c)
}

// Close graceful shutdown database
func (mdb *MultiDB) Close() {
	if mdb.aofHandler != nil {
		mdb.aofHandler.Close()
	}
}

func execSelect(c redis.Connection, mdb *MultiDB, args [][]byte) redis.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(mdb.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}

func (mdb *MultiDB) flushAll() redis.Reply {
	for _, db := range mdb.dbSet {
		db.Flush()
	}
	if mdb.aofHandler != nil {
		mdb.aofHandler.AddAof(0, utils.ToCmdLine("FlushAll"))
	}
	return &reply.OkReply{}
}

// ForEach traverses all the keys in the given database
func (mdb *MultiDB) ForEach(dbIndex int, cb func(key string, data *redis.DataEntity, expiration *time.Time) bool) {
	if dbIndex >= len(mdb.dbSet) {
		return
	}
	db := mdb.dbSet[dbIndex]
	db.ForEach(cb)
}

// ExecMulti executes multi commands transaction Atomically and Isolated
func (mdb *MultiDB) ExecMulti(conn redis.Connection, watching map[string]uint32, cmdLines [][][]byte) redis.Reply {
	if conn.GetDBIndex() >= len(mdb.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	db := mdb.dbSet[conn.GetDBIndex()]
	return db.ExecMulti(conn, watching, cmdLines)
}

// RWLocks lock keys for writing and reading
func (mdb *MultiDB) RWLocks(dbIndex int, writeKeys []string, readKeys []string) {
	if dbIndex >= len(mdb.dbSet) {
		panic("ERR DB index is out of range")
	}
	db := mdb.dbSet[dbIndex]
	db.RWLocks(writeKeys, readKeys)
}

// RWUnLocks unlock keys for writing and reading
func (mdb *MultiDB) RWUnLocks(dbIndex int, writeKeys []string, readKeys []string) {
	if dbIndex >= len(mdb.dbSet) {
		panic("ERR DB index is out of range")
	}
	db := mdb.dbSet[dbIndex]
	db.RWUnLocks(writeKeys, readKeys)
}

// GetUndoLogs return rollback commands
func (mdb *MultiDB) GetUndoLogs(dbIndex int, cmdLine [][]byte) [][][]byte {
	if dbIndex >= len(mdb.dbSet) {
		panic("ERR DB index is out of range")
	}
	db := mdb.dbSet[dbIndex]
	return db.GetUndoLogs(cmdLine)
}

// ExecWithLock executes normal commands, invoker should provide locks
func (mdb *MultiDB) ExecWithLock(conn redis.Connection, cmdLine [][]byte) redis.Reply {
	if conn.GetDBIndex() >= len(mdb.dbSet) {
		panic("ERR DB index is out of range")
	}
	db := mdb.dbSet[conn.GetDBIndex()]
	return db.execWithLock(cmdLine)
}

// BGRewriteAOF asynchronously rewrites Append-Only-File
func BGRewriteAOF(db *MultiDB, args [][]byte) redis.Reply {
	go db.aofHandler.Rewrite()
	return reply.MakeStatusReply("Background append only file rewriting started")
}

// RewriteAOF start Append-Only-File rewriting and blocked until it finished
func RewriteAOF(db *MultiDB, args [][]byte) redis.Reply {
	db.aofHandler.Rewrite()
	return reply.MakeStatusReply("Background append only file rewriting started")
}
