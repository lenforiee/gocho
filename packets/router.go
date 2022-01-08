package packets

import (
	"fmt"
	"strings"

	"github.com/lenforiee/gocho/constants"
	helpers "github.com/lenforiee/gocho/helpers"
	"github.com/valyala/fasthttp"
)

func LoginHandler(c *fasthttp.RequestCtx) (string, []byte) {
	// For debugging purposes for now.
	packet := BuildPacket(constants.BanchoLoginReply, int32(1))
	packet2 := BuildPacket(constants.BanchoAnnounce, "Welcome to gocho!")
	resp := append(packet, packet2...)

	data := strings.Split(string(c.Request.Body()), "\n")
	username := data[0]
	usernameSafe := helpers.UsernameSafe(username)
	password := data[1]
	moreData := strings.Split(data[2], "|")
	fmt.Println(data, username, usernameSafe, password, moreData)

	return "3h98r23823hr", resp
}

func RouterPage(c *fasthttp.RequestCtx) {
	if string(c.UserAgent()) != "osu!" {
		c.WriteString("gocho - an superior bancho written in go!")
		return
	}

	fmt.Println("Got da request.")
	choToken := string(c.Request.Header.Peek("osu-token"))
	if choToken == "" {
		uuid, packets := LoginHandler(c)
		c.Response.Header.Add("cho-token", uuid)
		c.Write(packets)
		return
	}

	body := c.Request.Body()
	var packets []Packet
	for {
		if len(body) == 0 {
			break
		}
		packets = append(packets, NewPacketReader(body))
		body = body[packets[len(packets)-1].PacketLen+7:]
	}
	fmt.Println(packets)
}
