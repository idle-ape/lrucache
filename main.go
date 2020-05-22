package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/idle-ape/lrucache/cache"
)

var (
	wg       sync.WaitGroup
	lruCache *cache.Cache
)

func init() {
	lruCache = cache.NewCache(1000)
}

func main() {
	for i := 1; i < 1001; i++ {
		wg.Add(1)
		go lruCacheTest(i)
	}

	wg.Wait()

	fmt.Println(len(lruCache.Cache), "prev", lruCache.Cache[strconv.Itoa(10)].Prev, "next", lruCache.Cache[strconv.Itoa(10)].Next, "head", lruCache.Head, "tail", lruCache.Tail)

	lruCache.Get(strconv.Itoa(10))

	fmt.Println(len(lruCache.Cache), "prev", lruCache.Cache[strconv.Itoa(10)].Prev, "next", lruCache.Cache[strconv.Itoa(10)].Next, "head", lruCache.Head, "tail", lruCache.Tail)

	lruCache.Set(strconv.Itoa(1200), make(chan string, 10), 0)

	fmt.Println(len(lruCache.Cache), "prev", lruCache.Cache[strconv.Itoa(1200)].Prev, "next", lruCache.Cache[strconv.Itoa(1200)].Next, "head", lruCache.Head, "tail", lruCache.Tail)
}

func lruCacheTest(i int) {

	rand.Seed(time.Now().Unix())

	lruCache.Set(strconv.Itoa(i), i, 0)

	lruCache.Get(strconv.Itoa(i))

	time.Sleep(1 * 1000 * time.Millisecond)

	lruCache.Get(strconv.Itoa(i))

	wg.Done()
}
