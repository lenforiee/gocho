package packets

import (
	//"fmt"
	"bytes"
	"encoding/binary"
)

type Packet struct {
	packetBuffer []byte
	offset int
	PacketID uint16
	PacketLen uint32
}

var HEADER_LEN = 7 

// Create new packet.
func NewPacketReader(buffer []byte) (packet Packet) {
	packet.packetBuffer = buffer
	packet.offset = 0
	packet.PacketID = 0
	packet.PacketLen = 0

	return packet
}

func ReadInt(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x) // cast type so we can return it.
	case 2:
		var x int16
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 4:
		var x int32
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	default:
		var x int64
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	}
}

func ReadUint(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x uint8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x) // cast type so we can return it.
	case 2:
		var x uint16
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 4:
		var x uint32
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	default:
		var x uint64
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	}
}

func (p *Packet) Read(size int) (data []byte) {
	data = p.packetBuffer[p.offset:p.offset+size]
	p.offset += size
	return data
}

func (p *Packet) ReadHeader() {
	//headerData := p.Read(HEADER_LEN)
	p.PacketID = p.ReadUint16()//int16(ReadInt(headerData[:2]))
	p.Read(1)
	p.PacketLen = p.ReadUint32()//int32(ReadInt(headerData[2:6]))
}

// func (p *Packet) MultiPacketCheck() {
// 	// Checks if the packet has additional data.
// 	if len(p.packetBuffer) > (p.packetLen + HEADER_LEN) {
// 		p.packetBuffer = p.packetBuffer[:HEADER_LEN + p.packetLen]
// 	}
// }

func (p *Packet) ReadUint8() uint8 {
	return uint8(ReadUint(p.Read(1))) // uncast the type.
}

func (p *Packet) ReadInt8() int8 {
	return int8(ReadInt(p.Read(1)))
}

func (p *Packet) ReadI32List() (vars []int32) {
	length := p.ReadUint16()
	for i := 0; i < int(length); i++ {
		vars = append(vars, p.ReadInt32())
	}
	return vars
}

func (p *Packet) ReadUint16() uint16 {
	return uint16(ReadUint(p.Read(2)))
}

func (p *Packet) ReadInt16() int16 {
	return int16(ReadInt(p.Read(2)))
}

func (p *Packet) ReadUint32() uint32 {
	return uint32(ReadUint(p.Read(4)))
}

func (p *Packet) ReadInt32() int32 {
	return int32(ReadInt(p.Read(4)))
}

func (p *Packet) ReadUint64() uint64 {
	return uint64(ReadUint(p.Read(8)))
}

func (p *Packet) ReadInt64() int64 {
	return int64(ReadInt(p.Read(8)))
}

func (p *Packet) ReadFloat64() (val float64) {
	buf := bytes.NewBuffer(p.Read(8))
	binary.Read(buf, binary.LittleEndian, &val)
	return val
}

func (p *Packet) ReadFloat32() (val float32) {
    buf := bytes.NewBuffer(p.Read(4))
	binary.Read(buf, binary.LittleEndian, &val)
	return val
}

func (p *Packet) ReadOsuString() string {
	if p.ReadUint8() != 0x0B {
		return ""
	}

	var val int
	var shift uint
	for {
		b := p.ReadUint8()
		val |= (int(b & 0x7F) << shift)
		if b & 0x80 == 0 {
			break
		}
		shift += 7
	}

	return string(p.Read(val))
}