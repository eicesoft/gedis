package database

import (
	"gedis/pkg/logger"
	"github.com/syndtr/goleveldb/leveldb"
)

type EachKey func([]byte, []byte)

type LevelDb struct {
	db *leveldb.DB
}

func NewLevelDb() *LevelDb {
	return &LevelDb{}
}

func (ldb *LevelDb) Load(path string) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		logger.Error(err.Error())
	}
	ldb.db = db
}

func (ldb *LevelDb) Put(key, value []byte) error {
	return ldb.db.Put(key, value, nil)
}

func (ldb *LevelDb) EachKeys(fn EachKey) {
	iter := ldb.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fn(key, value)
	}
	iter.Release()
}

func (ldb *LevelDb) Close() {
	ldb.db.Close()
}
