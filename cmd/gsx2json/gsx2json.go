package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"com.baby543.gsx2json-go/pkg/cache"
	"com.baby543.gsx2json-go/pkg/gsx2json"
	"com.baby543.gsx2json-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var opts struct {
	Port string `long:"port" short:"p" description:"Service port" default:"8080"`

	LogLevel int8 `long:"log-level" description:"referencing to zapcore.Level" default:"0"`

	CacheMode string `long:"cache" choice:"file" choice:"memory" choice:"none" default:"none"`

	SSLMode bool `long:"ssl" description:"referencing to zapcore.Level"`

	Version func() `long:"version" shot:"v" description:"display binary build version"`
}

var (
	Version = "0.0.0"
	Build   = "-"
)

var logger *zap.Logger
var caches cache.Manager
var parser = flags.NewParser(&opts, flags.Default)

func init() {
	opts.Version = func() {
		fmt.Printf("Version: %v", Version)
		fmt.Printf("\tBuild: %v", Build)
		os.Exit(0)
	}
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	switch opts.CacheMode {
	case "none":
		caches = cache.NewDummyCache()
	case "file":
		caches = cache.NewFileCache()
	case "memory":
		caches = cache.NewMemoryCache()
	}
	gin.SetMode(gin.ReleaseMode)
	utils.LogLevel = zapcore.Level(opts.LogLevel)
	logger = utils.NewLogger()
	defer logger.Sync()
}

func main() {
	log.Printf("Version: %v Build: %v", Version, Build)
	log.Printf("Port: %v", opts.Port)
	log.Printf("SSLMode: %v", opts.SSLMode)
	log.Printf("CacheMode: %v", opts.CacheMode)
	router := gin.Default()
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": Version,
			"build":   Build,
		})
	})
	router.DELETE("/cache", func(c *gin.Context) {
		defer caches.Flush()
		c.IndentedJSON(http.StatusOK, gin.H{
			"mode": opts.CacheMode,
			"size": len(caches.List()),
		})
	})
	router.GET("/cache", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"mode": opts.CacheMode,
			"size": len(caches.List()),
		})
	})
	router.GET("/api", func(c *gin.Context) {
		config := gsx2json.NewConfig()
		if err := config.Parse(c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
		identifier := gsx2json.NewIdentifier()
		if err := identifier.Parse(c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
		var err error
		var data []byte
		if data, err = caches.Load(identifier.String()); err != nil {
			data, err = gsx2json.Request(identifier)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err,
				})
				return
			}
			caches.Save(data, identifier.String())
		}
		payload := gsx2json.NewPayload()
		if err := payload.Parse(data, config); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err,
			})
			return
		}
		if config.PrettyPrint {
			c.IndentedJSON(http.StatusOK, payload.View)
		} else {
			c.JSON(http.StatusOK, payload.View)
		}
	})
	var err error
	addr := ":" + opts.Port
	if opts.SSLMode {
		err = router.RunTLS(addr,
			"./cert/cert.pem",
			"./cert/key.pem",
		)
	} else {
		err = router.Run(addr)
	}
	if err != nil {
		log.Fatal(err)
	}
}
