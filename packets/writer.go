package packets

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"reflect"
)

func WriteUleb128(value int) (ret []byte) {
	var len int

	if value == 0 {
		ret = []byte{0}
		return ret
	}

	for value > 0 {
		ret = append(ret, 0) // Stupid hack.
		ret[len] = byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			ret[len] |= 0x80
		}
		len++
	}
	return ret
}

func WriteOsuString(value string) (ret []byte) {

	if value == "" {
		ret = []byte{0}
	} else {
		b := []byte(value)
		ret = append(ret, 0x0B) // This is cursed.
		ret = append(ret, WriteUleb128(len(b))...)
		ret = append(ret, b...)
	}

	return ret
}

func BuildPacket(packetID int, contents ...interface{}) []byte {

	var packData []byte
	for _, content := range contents {
		switch typ := content.(type) {
		case *string, *[]byte, *[]int32, *[]uint32:
			BuildPacket(packetID, reflect.TypeOf(typ).Elem())
		case string:
			packData = append(packData, WriteOsuString(fmt.Sprint(content))...)
		}
	}

	packetLen := len(packData)
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, packetID)
	binary.Write(b, binary.LittleEndian, int8(0))
	binary.Write(b, binary.LittleEndian, packetLen)
	binary.Write(b, binary.LittleEndian, packData)

	return b.Bytes()
}
