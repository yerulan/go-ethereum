package mceliece

func store_gf(dest []byte, a uint16) {
	dest[0] = byte(a & 0xFF)
	dest[1] = byte(a >> 8)
}

func load_gf(src []byte) uint16 {
	a := uint16(src[1]) << 8
	a |= uint16(src[0])

	return a & GFMASK
}

func load4(in []byte) uint32 {
	ret := uint32(in[3])
	for i := 2; i >= 0; i-- {
		ret <<= 8
		ret |= uint32(in[i])
	}

	return ret
}

func store8(out []byte, in uint64) {
	out[0] = byte(in >> 0x00 & 0xFF)
	out[1] = byte(in >> 0x08 & 0xFF)
	out[2] = byte(in >> 0x10 & 0xFF)
	out[3] = byte(in >> 0x18 & 0xFF)
	out[4] = byte(in >> 0x20 & 0xFF)
	out[5] = byte(in >> 0x28 & 0xFF)
	out[6] = byte(in >> 0x30 & 0xFF)
	out[7] = byte(in >> 0x38 & 0xFF)
}

func load8(in []byte) uint64 {
	ret := uint64(in[7])
	for i := 6; i >= 0; i-- {
		ret <<= 8
		ret |= uint64(in[i])
	}

	return ret
}

func bitrev(a gf) gf {
	a = ((a & 0x00FF) << 8) | ((a & 0xFF00) >> 8)
	a = ((a & 0x0F0F) << 4) | ((a & 0xF0F0) >> 4)
	a = ((a & 0x3333) << 2) | ((a & 0xCCCC) >> 2)
	a = ((a & 0x5555) << 1) | ((a & 0xAAAA) >> 1)

	return a >> 3
}
