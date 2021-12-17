package database

import "strings"

//全局命令表
var cmdTable = make(map[string]*command)

type command struct {
	executor ExecFunc
	prepare  PreFunc // return related keys command
	undo     UndoFunc
	arity    int // allow number of args, arity < 0 means len(args) >= -arity
	flags    int
}

// RegisterCommand registers a new command
func RegisterCommand(name string, executor ExecFunc, prepare PreFunc, rollback UndoFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		prepare:  prepare,
		undo:     rollback,
		arity:    arity,
	}
}
