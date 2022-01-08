package packets

import (
	"bytes"
	"math"
	"encoding/binary"
)

type Packet struct {
	packetBuffer []byte
	offset int
	packetID int
	packetLen int
}

var HEADER_LEN = 7 

// Create new packet.
func NewPacketReader(buffer []byte) (packet Packet) {
	packet.packetBuffer = buffer
	packet.offset = 0
	packet.packetID = 0
	packet.packetLen = 0

	return packet
}

func ReadInt(b []byte) int {
	buf := bytes.NewBuffer(b)

	// I wanna hide this code somewhere so 
	// i dont have to look at it...
	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
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

	// I wanna hide this code somewhere so 
	// i dont have to look at it...
	switch len(b) {
	case 1:
		var x uint8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
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
	headerData := p.Read(HEADER_LEN)
	p.packetID = ReadInt(headerData[:2])
	p.packetLen = ReadInt(headerData[:4])
}

func (p *Packet) MultiPacketCheck() {
	// Checks if the packet has additional data.
	if len(p.packetBuffer) > (p.packetLen + HEADER_LEN) {
		p.packetBuffer = p.packetBuffer[:HEADER_LEN + p.packetLen]
	}
}

func (p *Packet) ReadUint8() int {
	return ReadUint(p.Read(1))
}

func (p *Packet) ReadInt8() int {
	return ReadInt(p.Read(1))
}

func (p *Packet) ReadI32List() (vars []int) {
	length := p.ReadUint16()
	for i := 0; i < length; i++ {
		vars = append(vars, p.ReadInt32())
	}
	return vars
}

func (p *Packet) ReadUint16() int {
	return ReadUint(p.Read(2))
}

func (p *Packet) ReadInt16() int {
	return ReadInt(p.Read(2))
}

func (p *Packet) ReadUint32() int {
	return ReadUint(p.Read(4))
}

func (p *Packet) ReadInt32() int {
	return ReadInt(p.Read(4))
}

func (p *Packet) ReadUint64() int {
	return ReadUint(p.Read(8))
}

func (p *Packet) ReadInt64() int {
	return ReadInt(p.Read(8))
}

func (p *Packet) ReadFloat64() float64 {
    bits := binary.LittleEndian.Uint64(p.Read(8))
    float := math.Float64frombits(bits)
    return float
}

func (p *Packet) ReadFloat32() float32 {
    bits := binary.LittleEndian.Uint32(p.Read(4))
    float := math.Float32frombits(bits)
    return float
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