package main

import (
	"fmt"
  "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/lenforiee/gocho/logger"
  packets "github.com/lenforiee/gocho/packets"
  constants "github.com/lenforiee/gocho/constants"
)

func rootPage(c *fasthttp.RequestCtx) {
	c.WriteString("gocho - an superior bancho written in go!")
}

func main() {
  router := router.New()

  router.GET("/", rootPage)

  logger.Info(fmt.Sprintf("Starting to handle gocho connections on 127.0.0.1:%d!", 8080))
  packet := packets.BuildPacket(constants.BanchoAnnounce, "Hello World", []int32{3, 10, 50})
  fmt.Println(packet)
  reader := packets.NewPacketReader(packet)
  reader.ReadHeader()
  fmt.Println(reader.PacketID)
  fmt.Println(reader.PacketLen)
  fmt.Println(reader.ReadOsuString())
  fmt.Println(reader.ReadI32List())
  fasthttp.ListenAndServe(":8080", router.Handler)
}