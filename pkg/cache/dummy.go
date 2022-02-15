package cache

import "errors"

type dummyCache struct {
	list []string
}

func NewDummyCache() Manager {
	return &dummyCache{
		list: make([]string, 0),
	}
}

func (c *dummyCache) Save(json []byte, key string) error {
	return errors.New("Cach mode disabled")
}

func (c *dummyCache) Load(key string) ([]byte, error) {
	return []byte{}, errors.New("Cach mode disabled")
}

func (c *dummyCache) List() []string {
	return c.list
}

func (c *dummyCache) Flush() {
}
