package packets

import (
	"fmt"
	bin "github.com/roman-kachanovsky/go-binary-pack/binary-pack"
)

type Packet struct {
	packetBuffer []byte
	offset int
	packetID int
	packetLen int
}

var HEADER_LEN = 7 

// Create new packet.
func newPacket(buffer []byte) (packet Packet) {
	packet.packetBuffer = buffer
	packet.offset = 0
	packet.packetID = 0
	packet.packetLen = 0

	return packet
}

func (p Packet) Read(size int) (data []byte) {
	data = p.packetBuffer[p.offset:p.offset+size]
	p.offset = (p.offset + size)
	return data
}

func (p Packet) ReadHeader() {
	headerData := p.Read(HEADER_LEN)
	bp := new(bin.BinaryPack)
	unpack, err := bp.UnPack([]string{"H", "I"}, headerData)
	if err != nil {
		return
	}
	fmt.Println(unpack)
}