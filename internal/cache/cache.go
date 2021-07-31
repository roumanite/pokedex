package cache

import (
	"time"
)

type Item struct {
	Object interface{}
	Expiration int64
}

type Cache struct {
	defaultExpiration time.Duration
	items map[string]Item
}

const (
	NoExpiration time.Duration = -1
	DefaultExpiration time.Duration = 0
)

func New(de time.Duration) *Cache {
	return &Cache{de, make(map[string]Item, 0)}
}

func (c *Cache) LoadFile(fname string) error {
	return nil
}

func (c *Cache) Write(fname string) error {
	return nil
} 

func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	return item.Object, true
}
