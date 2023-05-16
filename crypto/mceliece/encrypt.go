package mceliece

import (
	"unsafe"
)

func uint16_is_smaller_declassify(t uint16, u uint16) uint16 {
	mask := crypto_uint16_smaller_mask(crypto_uint16(t), crypto_uint16(u))
	crypto_declassify(unsafe.Pointer(&mask), unsafe.Sizeof(mask))
	return uint16(mask)
}

func uint32_is_equal_declassify(t uint32, u uint32) uint32 {
	mask := crypto_uint32_equal_mask(crypto_uint32(t), crypto_uint32(u))
	crypto_declassify(unsafe.Pointer(&mask), unsafe.Sizeof(mask))
	return uint32(mask)
}

func same_mask(x uint16, y uint16) uint8 {
	var mask uint32
	mask = uint32(x ^ y)
	mask -= 1
	mask >>= 31
	mask = -mask
	return uint8(mask & 0xFF)
}

func gen_e(e []byte) {
	var i, j, eq, count int

	type bufUnion struct {
		nums  [SYS_T * 2]uint16
		bytes [SYS_T * 2 * 2]byte
	}

	var buf bufUnion

	var ind [SYS_T]uint16
	var mask uint8
	var val [SYS_T]uint8

	for {
		randombytes(buf.bytes[:], uint64(unsafe.Sizeof(buf)))
		for i = 0; i < SYS_T*2; i++ {
			buf.nums[i] = load_gf(buf.bytes[i*2 : i*2+2])
		}

		count = 0
		for i = 0; i < SYS_T*2 && count < SYS_T; i++ {
			if uint16_is_smaller_declassify(buf.nums[i], SYS_N) != 0 {
				ind[count] = buf.nums[i]
				count++
			}
		}

		if count >= SYS_T {
			break
			continue
		}

		eq = 0
		for i = 1; i < SYS_T; i++ {
			for j = 0; j < i; j++ {
				if uint32_is_equal_declassify(uint32(ind[i]), uint32(ind[j])) != 0 {
					eq = 1
				}
			}
		}

		if eq == 0 {
			break
		}
	}

	for j = 0; j < SYS_T; j++ {
		val[j] = 1 << (ind[j] & 7)
	}

	for i = 0; i < SYS_N/8; i++ {
		e[i] = 0
		for j = 0; j < SYS_T; j++ {
			mask = same_mask(uint16(i), ind[j]>>3)
			e[i] |= val[j] & mask
		}
	}
}

func syndrome(s, pk, e []byte) {
	var b uint8
	row := make([]byte, SYS_N/8)
	pkPtr := 0

	for i := 0; i < SYND_BYTES; i++ {
		s[i] = 0
	}

	for i := 0; i < PK_NROWS; i++ {
		for j := 0; j < SYS_N/8; j++ {
			row[j] = 0
		}

		for j := 0; j < PK_ROW_BYTES; j++ {
			row[SYS_N/8-PK_ROW_BYTES+j] = pk[pkPtr]
			pkPtr++
		}

		row[i/8] |= 1 << (i % 8)

		b = 0
		for j := 0; j < SYS_N/8; j++ {
			b ^= row[j] & e[j]
		}
		b ^= b >> 4
		b ^= b >> 2
		b ^= b >> 1
		b &= 1

		s[i/8] |= (b << (i % 8))
	}
}
func encrypt(s, pk, e []byte) {
	gen_e(e)

	syndrome(s, pk, e)
}
