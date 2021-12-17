package database

import (
	"gedis/config"
	"gedis/reply"
	"gedis/types/cmd"
	"gedis/types/redis"
)

// Ping the server
func Ping(db *DB, args [][]byte) redis.Reply {
	if len(args) == 0 {
		return &reply.PongReply{}
	} else if len(args) == 1 {
		return reply.MakeStatusReply(string(args[0]))
	} else {
		return reply.MakeErrReply("ERR wrong number of arguments for 'ping' command")
	}
}

// Info 服务器信息输出
func Info(db *DB, args [][]byte) redis.Reply {
	return &reply.OkReply{}
}

// Auth validate client's password
func Auth(c redis.Connection, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'auth' command")
	}
	if config.Get().Server.Password == "" {
		return reply.MakeErrReply("ERR Client sent AUTH, but no password is set")
	}
	passwd := string(args[0])
	c.SetPassword(passwd)
	if config.Get().Server.Password != passwd {
		return reply.MakeErrReply("ERR invalid password")
	}
	return &reply.OkReply{}
}

func isAuthenticated(c redis.Connection) bool {
	if config.Get().Server.Password == "" {
		return true
	}
	return c.GetPassword() == config.Get().Server.Password
}

func init() {
	RegisterCommand(cmd.Ping, Ping, noPrepare, nil, -1)
	RegisterCommand(cmd.Info, Info, noPrepare, nil, -1)
}
