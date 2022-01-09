package packets

import (
	"fmt"
	"strings"

	"github.com/lenforiee/gocho/constants"
	"github.com/lenforiee/gocho/helpers"
	"github.com/lenforiee/gocho/structures"
	"github.com/valyala/fasthttp"
)

type basicInfo struct {
	UserID   int    `db:"id"`
	Md5Crypt string `db:"password_md5"`
}

func LoginHandler(c *fasthttp.RequestCtx) (string, []byte) {
	// For debugging purposes for now.
	// packet := BuildPacket(constants.BanchoLoginReply, int32(-1))
	// packet2 := BuildPacket(constants.BanchoAnnounce, "Welcome to gocho!")
	// resp := append(packet, packet2...)
	//var resp []byte
	data := strings.Split(string(c.Request.Body()), "\n")
	username := data[0]
	usernameSafe := helpers.UsernameSafe(username)
	password := data[1]
	moreData := strings.Split(data[2], "|")
	fmt.Println(data, username, usernameSafe, password, moreData)

	// info := basicInfo{}
	// err := globals.DbConn.Get(&info, "SELECT id, password_md5 FROM users WHERE username_safe = ? LIMIT 1", usernameSafe)
	// switch {
	// case err == sql.ErrNoRows: // TODO: Write normal packets contexts.
	// 	resp = append(resp, BuildPacket(constants.BanchoLoginReply, int32(-1))...)
	// 	resp = append(resp, BuildPacket(constants.BanchoAnnounce, "gocho: This username doesn't exist!")...)
	// 	return "no", resp
	// case err != nil:
	// 	resp = append(resp, BuildPacket(constants.BanchoLoginReply, int32(-5))...)
	// 	return "no", resp
	// }

	// They exists!
	// err = bcrypt.CompareHashAndPassword([]byte(info.Md5Crypt), []byte(password))
	// if err != nil {
	// 	resp = append(resp, BuildPacket(constants.BanchoLoginReply, int32(-1))...)
	// 	resp = append(resp, BuildPacket(constants.BanchoAnnounce, "gocho: Wrong Password!")...)
	// 	return "no", resp
	// }
	u, _ := structures.NewUser(1)
	u.Queue(BuildPacket(constants.BanchoLoginReply, int32(1)))
	u.Queue(BuildPacket(constants.BanchoProtocolNegotiation, int32(19)))
	u.Queue([]byte("Y\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00"))
	u.Queue(BuildPacket(constants.BanchoBanInfo, uint32(0)))
	u.Queue(BuildPacket(constants.BanchoLoginPermissions, uint32(5)))

	presense := BuildPacket(constants.BanchoUserPresence, int32(1), "lenforiee", uint8(25), uint8(2), uint8(5), float32(39.01955903386848), float32(125.75276158057767), int32(0))
	presense = append(presense, BuildPacket(constants.BanchoHandleOsuUpdate, int32(1), uint(0), "", "", "", int32(0), uint8(0), int32(0), int64(0), float32(0.0), int32(0), int64(0), int32(0), int16(0))...)

	structures.Broadcast(presense)
	u.Queue(presense)
	u.Queue(BuildPacket(constants.BanchoAnnounce, "gocho! login successful!"))

	return u.UUID, u.Dequeue()
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
