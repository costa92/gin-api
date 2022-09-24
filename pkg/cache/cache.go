package cache

import (
	"fmt"
	"sync"

	"github.com/costa92/errors"
	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	lock    *sync.RWMutex
	secrets *ristretto.Cache
}

// ErrSecretNotFound defines secret not found error.
var ErrSecretNotFound = errors.New("secret not found")

var (
	onceCache sync.Once
	cacheIns  *Cache
)

func GetCacheInsOr() (*Cache, error) {
	var err error
	var secretCache *ristretto.Cache
	onceCache.Do(func() {
		c := &ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
			Cost:        nil,
		}

		secretCache, err = ristretto.NewCache(c)
		if err != nil {
			return
		}
		cacheIns = &Cache{
			lock:    new(sync.RWMutex),
			secrets: secretCache,
		}
	})
	return cacheIns, err
}

func (c *Cache) GetSecret(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.secrets.Get(key)
	if !ok {
		fmt.Println(ok)
		return nil
	}
	fmt.Println(value)
	return nil
}

func (c *Cache) SetSecret(key, val string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.secrets.Clear()
	c.secrets.Set(key, val, 1)
	c.secrets.Wait()
}
