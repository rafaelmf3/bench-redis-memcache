package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/rainycape/memcache"
)

var num = 10000

func BenchmarkRedis(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := client.Set("key", "data", 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	var wg *sync.WaitGroup
	wg = &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		redisWorker(client, wg, num)
	}
	wg.Wait()
}

func BenchmarkMemcached(b *testing.B) {
	memcachedClient, _ := memcache.New("127.0.0.1:11211")
	memcachedClient.SetTimeout(time.Second * 10)
	memcachedClient.SetMaxIdleConnsPerAddr(100)
	memcachedClient.Set(&memcache.Item{Key: "key", Value: []byte("data")})
	var wg *sync.WaitGroup
	wg = &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		memcachedWorker(memcachedClient, wg, num)
	}
	wg.Wait()
}
