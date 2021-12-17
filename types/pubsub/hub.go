package pubsub

import (
	"gedis/types/dict"
	"gedis/types/lock"
)

type Hub struct {
	subs dict.Dict // channel -> list(*Client)

	subsLocker *lock.Locks // lock channel

}

// MakeHub creates new hub
func MakeHub() *Hub {
	return &Hub{
		subs:       dict.MakeConcurrent(4),
		subsLocker: lock.Make(16),
	}
}
