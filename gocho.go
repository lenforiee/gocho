package main

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/lenforiee/gocho/logger"
	packets "github.com/lenforiee/gocho/packets"
)

func RootPage(c *fasthttp.RequestCtx) {
	c.WriteString("gocho - an superior bancho written in go!")
}

func main() {
	router := router.New()

	router.GET("/", RootPage)
	router.POST("/", packets.RouterPage)

	logger.Info(fmt.Sprintf("Starting to handle gocho connections on 127.0.0.1:%d!", 8080))
	fasthttp.ListenAndServe(":8080", router.Handler)
}
