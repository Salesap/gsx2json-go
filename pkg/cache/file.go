package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const SAVE_DIR = ".cache/"

type fileCache struct {
	files map[string]bool
}

func init() {
	os.Mkdir(SAVE_DIR, os.ModePerm)
}

func NewFileCache() Manager {
	c := &fileCache{
		files: make(map[string]bool),
	}
	if files, err := ioutil.ReadDir(SAVE_DIR); err != nil {
		for _, info := range files {
			if filepath.Ext(info.Name()) == "json" {
				c.files[info.Name()] = true
			}
		}
	}
	return c
}

func (c *fileCache) Save(json []byte, key string) error {
	const flags = os.O_WRONLY | os.O_CREATE
	var filepath = SAVE_DIR + key + ".json"
	f, err := os.OpenFile(filepath, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	c.files[key] = true
	_, err = f.Write(json)
	return err
}

func (c *fileCache) Load(key string) ([]byte, error) {
	b := []byte{}
	const flags = os.O_RDONLY
	var filepath = SAVE_DIR + key + ".json"
	f, err := os.OpenFile(filepath, flags, 0644)
	if err != nil {
		return b, err
	}
	defer f.Close()
	b, err = ioutil.ReadAll(f)
	return b, err
}

func (c *fileCache) List() []string {
	list := make([]string, 0)
	for k, _ := range c.files {
		list = append(list, k)
	}
	return list
}

func (c *fileCache) Flush() {
	os.RemoveAll(SAVE_DIR)
	c.files = make(map[string]bool)
}
