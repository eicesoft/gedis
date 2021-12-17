package server

import (
	"context"
	"gedis/database"
	"gedis/pkg/utils"
	"gedis/types/redis"
	"gedis/types/redis/connection"
	"net"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

const (
	Version = "0.9.1"
)

type Server struct {
	addr       string        // 监听地址
	db         redis.DB      // DB结构
	closing    utils.Boolean //
	activeConn sync.Map      // *client -> placeholder
}

type IServer interface {
	Run()
	Handle(ctx context.Context, conn net.Conn)
	closeClient(client *connection.Connection)
}

func NewServer(addr string) *Server {
	var db redis.DB

	db = database.NewStandaloneServer() //单机模式

	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) closeClient(client *connection.Connection) {
	_ = client.Close()
	s.db.AfterClientClose(client)
	s.activeConn.Delete(client)
}

func (s *Server) Handle(ctx context.Context, conn net.Conn) {
	if s.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
	}

	client := connection.NewConn(conn)
	s.activeConn.Store(client, 1)
	client.ProcessCommand(func(buf *[][]byte) redis.Reply {
		return s.db.Exec(client, *buf)
	}, func() {
		s.closeClient(client)
	})
}

func (s *Server) Close() error {
	s.closing.Set(true)
	s.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	s.db.Close()

	return nil
}

// Run 服务运行
func (s *Server) Run() {
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			s.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
