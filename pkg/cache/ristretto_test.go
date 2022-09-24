package cache

import (
	"fmt"
	"testing"

	"github.com/dgraph-io/ristretto"
)

func Test_ristretto(t *testing.T) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	cache.Set("key", "value", 1)

	cache.Wait()
	value, found := cache.Get("key")
	if !found {
		panic("missing value")
	}
	fmt.Println(value)
	cache.Del("key")
}
