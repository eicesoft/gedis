package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		//Addr:         "127.0.0.1:6379",
		Addr:         "192.168.1.21:6379",
		Password:     "",
		DB:           0,
		MaxRetries:   1,
		PoolSize:     5,
		MinIdleConns: 2,
	})
	client.Set(context.Background(), "A1", "sdgasdgsadg", 30*time.Minute)
	d := client.Get(context.Background(), "A1")
	t.Logf("%v", d)
}

func TestInterface1(t *testing.T) {
	table := make(map[string]interface{})
	table["a1"] = 1
	x := table["a1"].(*int32)

	v := atomic.AddInt32(x, 1)
	fmt.Printf("%v, %d", table, v)
}

func BenchmarkClient(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:         "192.168.1.21:6381",
		Password:     "",
		DB:           0,
		MaxRetries:   1,
		PoolSize:     5,
		MinIdleConns: 2,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set(context.Background(), fmt.Sprintf("K:%d", i), fmt.Sprintf("vvv:%d", i), time.Second*0)
	}
}

func BenchmarkRedisClient(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:         "192.168.1.21:6379",
		Password:     "",
		DB:           0,
		MaxRetries:   1,
		PoolSize:     5,
		MinIdleConns: 2,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set(context.Background(), fmt.Sprintf("K:%d", i), fmt.Sprintf("vvv:%d", i), time.Second*0)
	}
}
