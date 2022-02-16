package gsx2json

import "github.com/gin-gonic/gin"

type Config struct {
	Query       string
	UseInteger  bool
	ShowDict    bool
	ShowRows    bool
	ShowColumns bool
	BriefMeta   bool
	PrettyPrint bool
}

func NewConfig() *Config {
	return &Config{
		Query:       "",
		UseInteger:  true,
		ShowDict:    true,
		ShowRows:    true,
		ShowColumns: true,
		BriefMeta:   false,
		PrettyPrint: false,
	}
}

func (cfg *Config) Parse(c *gin.Context) error {
	if v, ok := c.GetQuery("q"); ok {
		cfg.Query = v
	}
	if v, ok := c.GetQuery("integers"); ok {
		cfg.UseInteger = v == "true"
	}
	if v, ok := c.GetQuery("dict"); ok {
		cfg.ShowDict = v == "true"
	}
	if v, ok := c.GetQuery("rows"); ok {
		cfg.ShowRows = v == "true"
	}
	if v, ok := c.GetQuery("columns"); ok {
		cfg.ShowColumns = v == "true"
	}
	if v, ok := c.GetQuery("meta"); ok {
		cfg.BriefMeta = v == "true"
	}
	if v, ok := c.GetQuery("pretty"); ok {
		cfg.PrettyPrint = v == "true"
	}
	return nil
}
