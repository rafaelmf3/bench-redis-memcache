package main

import (
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/go-redis/redis"

	//	"runtime"

	"github.com/rainycape/memcache"
)

type ResourceConn struct {
	redis.Conn
}

func (r ResourceConn) Close() {
	r.Conn.Close()
}

func main() {
	//	runtime.GOMAXPROCS(1)

	parallel := 100
	iter := 100

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := client.Set("key", "data", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	var start, end time.Time
	var wg *sync.WaitGroup

	log.Println("Start redis test bench")
	wg = &sync.WaitGroup{}
	start = time.Now()
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go redisWorker(client, wg, iter)
	}
	wg.Wait()
	end = time.Now()
	log.Println("requests per second", int64(parallel*iter)/((end.UnixNano()-start.UnixNano())/1000000000))

	memcachedClient, err := memcache.New("127.0.0.1:11211")
	memcachedClient.SetTimeout(time.Second * 10)
	memcachedClient.SetMaxIdleConnsPerAddr(100)
	memcachedClient.Set(&memcache.Item{Key: "key", Value: []byte("data")})
	if err != nil {
		panic(err)
	}

	log.Println("Start memcached test bench")
	wg = &sync.WaitGroup{}
	start = time.Now()
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go memcachedWorker(memcachedClient, wg, iter)
	}
	wg.Wait()
	end = time.Now()
	log.Println("requests per second", int64(parallel*iter)/((end.UnixNano()-start.UnixNano())/1000000000))

}

func memcachedWorker(m *memcache.Client, wg *sync.WaitGroup, n int) {
	for i := 0; i < n; i++ {

		_, err := m.Get("key")

		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func redisWorker(c *redis.Client, wg *sync.WaitGroup, n int) {
	for i := 0; i < n; i++ {
		_, err := c.Get("key").Result()
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}
