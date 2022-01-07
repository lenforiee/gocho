package packets

func WriteUleb128(value int) (ret []byte) {
	var len int

	if value == 0 {
		ret = []byte{0}
		return ret
	}

	for value > 0 {
		ret = append(ret, 0) // Stupid hack.
		ret[len] = byte(value & 127)
		value >>= 7
		if value != 0 {
			ret[len] |= 128
		}
		len++
	}

	return ret
}

func WriteOsuString(value string) (ret []byte) {

	if value == "" {
		ret = []byte{0}
	} else {
		ret = append(ret, "\x0B"...) // This is cursed.
		ret = append(ret[:], WriteUleb128(len(value))[:]...)
		ret = append(ret, value...)
	}

	return ret
}
