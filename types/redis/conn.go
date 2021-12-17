package redis

import (
	"net"
)

type Func func(*[][]byte) Reply
type CloseFunc func()

type Connection interface {
	RemoteAddr() net.Addr

	Write([]byte) error
	SetPassword(string)
	GetPassword() string

	// client should keep its subscribing channels
	Subscribe(channel string)
	UnSubscribe(channel string)
	SubsCount() int
	GetChannels() []string

	InMultiState() bool
	SetMultiState(bool)
	GetQueuedCmdLine() [][][]byte
	EnqueueCmd([][]byte)
	ClearQueuedCmds()
	GetWatching() map[string]uint32
	GetDBIndex() int
	SelectDB(int)
	ProcessCommand(fn Func, closeFn CloseFunc)
}
