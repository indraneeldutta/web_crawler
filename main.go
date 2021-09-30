package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/web_crawler/apis"
)

func main() {

}

func applicationConfig() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")

	apis.NewApiController(v1)

	router.Run(":9000")
}
