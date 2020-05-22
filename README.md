# lrucache
Lrucache is an efficient cache library implemented by lru algorithm.

Usage:

```go
// import lru cache library
import(
    "fmt"
    "github.com/idle-ape/lrucache/cache"
)

// initial a fixed size cache
lrucache := cache.NewCache(1000)

lruCache.Set("key", "value", 0)

fmt.Println(lruCache.Get("key"))
```

