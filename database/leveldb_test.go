package database

import (
	"fmt"
	"testing"
	"time"
)

func TestLevelDB(t *testing.T) {
	start := time.Now()
	db := NewLevelDb()
	db.Load("./data")
	defer db.Close()

	//for i := 0; i < 1000000; i++ {
	//	key := fmt.Sprintf("a:%d", i)
	//	db.Put([]byte(key), []byte("hello world, hello worldhello worldhello worldhello worldhello worldhello world"))
	//}

	db.EachKeys(func(key []byte, val []byte) {
		fmt.Println(string(key))
	})

	fmt.Printf("Run time is %v", time.Since(start))
}
