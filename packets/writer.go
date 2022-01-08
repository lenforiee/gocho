package packets

import (
	//"fmt"
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
		ret = append(ret, 0x0B)
		ret = append(ret, WriteUleb128(len(b))...)
		ret = append(ret, b...)
	}

	return ret
}

func BuildPacket(packetID int, contents ...interface{}) []byte {

	packData := new(bytes.Buffer)
	for _, content := range contents {
		switch typ := content.(type) {
		case *string, *[]byte, *[]int32:
			BuildPacket(packetID, reflect.TypeOf(typ).Elem())
		case string, []byte:
			casted, _ := content.(string)
			packData.Write(WriteOsuString(casted))
		case []int32:
			casted, _ := content.([]int32)
			binary.Write(packData, binary.LittleEndian, uint16(len(casted)))
			for _, val := range casted {
				binary.Write(packData, binary.LittleEndian, val)
			}
		default:
			binary.Write(packData, binary.LittleEndian, content)
		}
	}

	b := new(bytes.Buffer)
	packBytes := packData.Bytes()
	binary.Write(b, binary.LittleEndian, uint16(packetID))
	binary.Write(b, binary.LittleEndian, int8(0))
	binary.Write(b, binary.LittleEndian, uint32(len(packBytes)))
	binary.Write(b, binary.LittleEndian, packBytes)

	return b.Bytes()
}
