package cache

import (
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"time"
	"os"
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
	DefaultExpiration time.Duration = -1
	NoExpiration time.Duration = 0
)

func New(de time.Duration) *Cache {
	return &Cache{de, make(map[string]Item, 0)}
}

func (c *Cache) LoadFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := scanner.Text()
		tokens := strings.Split(str, "|")

		if len(tokens) > 2 {
			duration, err := strconv.ParseInt(tokens[2], 10, 64)
			if err != nil {
				continue
			}
			c.items[tokens[0]] = Item{
				Object:     []byte(tokens[1]),
				Expiration: duration,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Write(fname string) error {
	f, err := os.Create(fname)
  if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	
	for k, v := range c.items {
		w.WriteString(fmt.Sprintf("%s|%s|%d\n", k, v.Object, v.Expiration))
	}

	w.Flush()
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
