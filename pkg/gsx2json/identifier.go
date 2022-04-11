package gsx2json

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
)

type Identifier struct {
	SheetId   string
	SheetName string
	ApiKey    string
}

func NewIdentifier() *Identifier {
	return &Identifier{}
}

func (id *Identifier) Parse(c *gin.Context) error {
	if v, ok := c.GetQuery("id"); ok {
		id.SheetId = v
	}
	if v, ok := c.GetQuery("sheet"); ok {
		id.SheetName = v
	}
	if v, ok := c.GetQuery("api_key"); ok {
		id.ApiKey = v
	}
	if len(id.SheetId) == 0 {
		return errors.New("You must provide a sheet ID.")
	}
	if len(id.SheetName) == 0 {
		return errors.New("You must provide a sheet name.")
	}
	if len(id.ApiKey) == 0 {
		if v, ok := os.LookupEnv("API_KEY"); ok {
			id.ApiKey = v
		}
		if len(id.ApiKey) == 0 {
			return errors.New("You must provide an api key.")
		}
	}
	return nil
}

func (id *Identifier) String() string {
	return id.SheetId + "_" + id.SheetName
}
