package redis

import "time"

// DB is the interface for redis style storage engine
type DB interface {
	Exec(client Connection, args [][]byte) Reply
	AfterClientClose(c Connection)
	Close()
}

type EmbedDB interface {
	DB
	ExecWithLock(conn Connection, args [][]byte) Reply
	ExecMulti(conn Connection, watching map[string]uint32, cmdLines [][][]byte) Reply
	GetUndoLogs(dbIndex int, cmdLine [][]byte) [][][]byte
	ForEach(dbIndex int, cb func(key string, data *DataEntity, expiration *time.Time) bool)
	RWLocks(dbIndex int, writeKeys []string, readKeys []string)
	RWUnLocks(dbIndex int, writeKeys []string, readKeys []string)
}

type DataEntity struct {
	Data interface{}
}
