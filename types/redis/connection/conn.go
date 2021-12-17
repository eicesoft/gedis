package connection

import (
	"bytes"
	"fmt"
	"gedis/pkg/logger"
	"gedis/reply"
	"gedis/server/parser"
	"gedis/types/redis"
	"io"
	"net"
	"strings"
	"sync"
)

const DefaultDbIndex = 0

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type Connection struct {
	conn       net.Conn
	mu         sync.Mutex
	password   string
	selectedDB int
	subs       map[string]bool // subscribing channels
	queue      [][][]byte
	watching   map[string]uint32
	multiState bool // queued commands for `multi`
	redis.Connection
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn:       conn,
		selectedDB: DefaultDbIndex,
		password:   "",
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) Write(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	c.mu.Lock()
	//c.waitingReply.Add(1)
	defer func() {
		//	c.waitingReply.Done()
		c.mu.Unlock()
	}()

	_, err := c.conn.Write(b)
	return err
}

// GetChannels returns all subscribing channels
func (c *Connection) GetChannels() []string {
	if c.subs == nil {
		return make([]string, 0)
	}
	channels := make([]string, len(c.subs))
	i := 0
	for channel := range c.subs {
		channels[i] = channel
		i++
	}
	return channels
}

// SetPassword set password for authentication
func (c *Connection) SetPassword(password string) {
	c.password = password
}

// GetPassword get password for authentication
func (c *Connection) GetPassword() string {
	return c.password
}

// InMultiState tells is connection in an uncommitted transaction
func (c *Connection) InMultiState() bool {
	return c.multiState
}

// SetMultiState sets transaction flag
func (c *Connection) SetMultiState(state bool) {
	if !state { // reset data when cancel multi
		c.watching = nil
		c.queue = nil
	}
	c.multiState = state
}

// GetQueuedCmdLine returns queued commands of current transaction
func (c *Connection) GetQueuedCmdLine() [][][]byte {
	return c.queue
}

// EnqueueCmd  enqueues command of current transaction
func (c *Connection) EnqueueCmd(cmdLine [][]byte) {
	c.queue = append(c.queue, cmdLine)
}

// ClearQueuedCmds clears queued commands of current transaction
func (c *Connection) ClearQueuedCmds() {
	c.queue = nil
}

// GetWatching returns watching keys and their version code when started watching
func (c *Connection) GetWatching() map[string]uint32 {
	if c.watching == nil {
		c.watching = make(map[string]uint32)
	}
	return c.watching
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

// SelectDB selects a database
func (c *Connection) SelectDB(dbNum int) {
	c.selectedDB = dbNum
}

// Subscribe add current connection into subscribers of the given channel
func (c *Connection) Subscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.subs == nil {
		c.subs = make(map[string]bool)
	}
	c.subs[channel] = true
}

// UnSubscribe removes current connection into subscribers of the given channel
func (c *Connection) UnSubscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.subs) == 0 {
		return
	}
	delete(c.subs, channel)
}

// SubsCount returns the number of subscribing channels
func (c *Connection) SubsCount() int {
	return len(c.subs)
}

func (c *Connection) ProcessCommand(fn redis.Func, closeFn redis.CloseFunc) {
	ch := parser.ParseStream(c.conn)

	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				closeFn() // connection closed
				logger.Info("connection closed: " + c.RemoteAddr().String())
				return
			}
			// protocol err
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := c.Write(errReply.ToBytes())
			if err != nil {
				closeFn() // connection closed
				logger.Info("connection closed: " + c.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			logger.Warn("empty payload")
			continue
		}
		r, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Warn("require multi bulk reply")
			continue
		}
		logger.Debug(fmt.Sprintf("执行命令: %s", r))
		result := fn(&r.Args)
		if result != nil {
			_ = c.Write(result.ToBytes())
		} else {
			_ = c.Write(unknownErrReplyBytes)
		}
	}
}

// FakeConn implements redis.Connection for test
type FakeConn struct {
	Connection
	buf bytes.Buffer
}

// Write writes data to buffer
func (c *FakeConn) Write(b []byte) error {
	c.buf.Write(b)
	return nil
}

// Clean resets the buffer
func (c *FakeConn) Clean() {
	c.buf.Reset()
}

// Bytes returns written data
func (c *FakeConn) Bytes() []byte {
	return c.buf.Bytes()
}
