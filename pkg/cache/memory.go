package cache

import (
	"bytes"
	"fmt"
)

type memoryCache struct {
	cache map[string]bytes.Buffer
}

func NewMemoryCache() Manager {
	return &memoryCache{
		cache: make(map[string]bytes.Buffer),
	}
}

func (c *memoryCache) Save(json []byte, key string) error {
	b := bytes.Buffer{}
	_, err := b.Read(json)
	if err != nil {
		return err
	}
	c.cache[key] = b
	return nil
}

func (c *memoryCache) Load(key string) ([]byte, error) {
	if b, ok := c.cache[key]; ok {
		return b.Bytes(), nil
	}
	return []byte{}, fmt.Errorf("%s not exist", key)
}

func (c *memoryCache) List() []string {
	list := make([]string, 0)
	for k, _ := range c.cache {
		list = append(list, k)
	}
	return list
}

func (c *memoryCache) Flush() {
	c.cache = make(map[string]bytes.Buffer)
}
