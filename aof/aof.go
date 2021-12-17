package aof

import (
	"gedis/config"
	"gedis/pkg/logger"
	"gedis/pkg/utils"
	"gedis/reply"
	"gedis/server/parser"
	"gedis/types/redis"
	"gedis/types/redis/connection"
	"io"
	"os"
	"strconv"
	"sync"
)

type CmdLine = [][]byte

const (
	aofQueueSize = 1 << 16
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// Handler 从aofChan通道中接受消息并写入文件处理器
type Handler struct {
	db          redis.EmbedDB //数据db, 恢复数据用
	tmpDBMaker  func() redis.EmbedDB
	aofChan     chan *payload //aof数据通道
	aofFile     *os.File      //aof数据文件
	aofFilename string        //aof数据文件名
	aofFinished chan struct{} //aof写入完成通道
	pausingAof  sync.RWMutex  //数据写入锁
	currentDB   int           //当前Db index
}

// NewAOFHandler new aof.Handler
func NewAOFHandler(db redis.EmbedDB, tmpDBMaker func() redis.EmbedDB) (*Handler, error) {
	handler := &Handler{}
	handler.aofFilename = "gedis.aof"
	handler.db = db
	handler.tmpDBMaker = tmpDBMaker
	handler.LoadAof(0)
	aofFile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	handler.aofChan = make(chan *payload, aofQueueSize)
	handler.aofFinished = make(chan struct{})
	go func() {
		handler.handleAof()
	}()
	return handler, nil
}

// AddAof 发生命令记录到通道中
func (handler *Handler) AddAof(dbIndex int, cmdLine CmdLine) {
	if config.Get().Server.EnableAof && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}
}

// handleAof 监听aof通道并写入文件
func (handler *Handler) handleAof() {
	// serialized execution
	handler.currentDB = 0
	for p := range handler.aofChan {
		handler.pausingAof.RLock() // prevent other goroutines from pausing aof
		if p.dbIndex != handler.currentDB {
			// select db
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("SELECT", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Warn(err.Error())
				continue // skip this command
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Warn(err.Error())
		}
		handler.pausingAof.RUnlock()
	}
	handler.aofFinished <- struct{}{}
}

// LoadAof 加载Aof文件
func (handler *Handler) LoadAof(maxBytes int) {
	logger.Debug("Load Aof file to memory start.")
	// delete aofChan to prevent write again
	aofChan := handler.aofChan
	handler.aofChan = nil
	defer func(aofChan chan *payload) {
		handler.aofChan = aofChan
	}(aofChan)

	file, err := os.Open(handler.aofFilename)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return
		}
		logger.Warn(err.Error())
		return
	}
	defer file.Close()

	var reader io.Reader
	if maxBytes > 0 {
		reader = io.LimitReader(file, int64(maxBytes))
	} else {
		reader = file
	}
	ch := parser.ParseStream(reader)
	fakeConn := &connection.FakeConn{} // only used for save dbIndex
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			logger.Warn("parse error: " + p.Err.Error())
			continue
		}
		if p.Data == nil {
			logger.Warn("empty payload")
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Warn("require multi bulk reply")
			continue
		}
		ret := handler.db.Exec(fakeConn, r.Args)
		if reply.IsErrorReply(ret) {
			logger.Warn("exec err: " + err.Error())
		}
	}
	logger.Debug("Load Aof file to memory success.")
}

// Close gracefully stops aof persistence procedure
func (handler *Handler) Close() {
	if handler.aofFile != nil {
		close(handler.aofChan)
		<-handler.aofFinished // wait for aof finished
		err := handler.aofFile.Close()
		if err != nil {
			logger.Warn(err.Error())
		}
	}
}
