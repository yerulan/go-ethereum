package mceliece

import "unsafe"

func layer_in(data [2][64]uint64, bits []uint64, lgs int) {
	s := 1 << lgs

	for i := 0; i < 64; i += s * 2 {
		for j := i; j < i+s; j++ {
			d := (data[0][j+0] ^ data[0][j+s])
			d &= bits[0]
			data[0][j+0] ^= d
			data[0][j+s] ^= d

			d = (data[1][j+0] ^ data[1][j+s])
			d &= bits[1]
			data[1][j+0] ^= d
			data[1][j+s] ^= d

			bits = bits[2:]
		}
	}
}

func layer_ex(data []uint64, bits []uint64, lgs int) {
	s := 1 << lgs

	for i := 0; i < 128; i += s * 2 {
		for j := i; j < i+s; j++ {
			d := (data[j+0] ^ data[j+s])
			d &= bits[0]
			data[j+0] ^= d
			data[j+s] ^= d

			bits = bits[1:]
		}
	}
}

func apply_benes(r, bits []byte, rev int) {
	inc := 0
	bitsPtr := bits

	if rev != 0 {
		bitsPtr = bits[12288:]
		inc = -1024
	}

	var (
		rIntV [2][64]uint64
		rIntH [2][64]uint64
		bIntV [64]uint64
		bIntH [64]uint64
	)

	rPtr := (*[1 << 30]byte)(unsafe.Pointer(&r[0]))

	for i := 0; i < 64; i++ {
		rIntV[0][i] = load8(rPtr[i*16+0:])
		rIntV[1][i] = load8(rPtr[i*16+8:])
	}

	transpose64x64(rIntH[0][:], rIntV[0][:])
	transpose64x64(rIntH[1][:], rIntV[1][:])

	for iter := 0; iter <= 6; iter++ {
		for i := 0; i < 64; i++ {
			bIntV[i] = load8(bitsPtr)
			bitsPtr = bitsPtr[8:]
		}

		bitsPtr = bitsPtr[inc:]

		transpose64x64(bIntH[:], bIntV[:])

		layer_ex(rIntH[0][:], bIntH[:], iter)
	}

	transpose64x64(rIntV[0][:], rIntH[0][:])
	transpose64x64(rIntV[1][:], rIntH[1][:])

	for iter := 0; iter <= 5; iter++ {
		for i := 0; i < 64; i++ {
			bIntV[i] = load8(bitsPtr)
			bitsPtr = bitsPtr[8:]
		}

		bitsPtr = bitsPtr[inc:]

		layer_in(rIntV, bIntV[:], iter)
	}

	for iter := 4; iter >= 0; iter-- {
		for i := 0; i < 64; i++ {
			bIntV[i] = load8(bitsPtr)
			bitsPtr = bitsPtr[8:]
		}
		bitsPtr = bitsPtr[inc:]

		layer_in(rIntV, bIntV[:], iter)
	}

	transpose64x64(rIntH[0][:], rIntV[0][:])
	transpose64x64(rIntH[1][:], rIntV[1][:])

	for iter := 6; iter >= 0; iter-- {
		for i := 0; i < 64; i++ {
			bIntV[i] = load8(bitsPtr)
			bitsPtr = bitsPtr[8:]
		}

		bitsPtr = bitsPtr[inc:]

		transpose64x64(bIntH[:], bIntV[:])

		layer_ex(rIntH[0][:], bIntH[:], iter)
	}

	transpose64x64(rIntV[0][:], rIntH[0][:])
	transpose64x64(rIntV[1][:], rIntH[1][:])

	for i := 0; i < 64; i++ {
		store8(rPtr[i*16+0:], rIntV[0][i])
		store8(rPtr[i*16+8:], rIntV[1][i])
	}
}

func support_gen(s []gf, c []byte) {
	var L [GFBITS][1 << GFBITS / 8]byte
	for i := 0; i < GFBITS; i++ {
		for j := 0; j < 1<<GFBITS/8; j++ {
			L[i][j] = 0
		}
	}

	for i := 0; i < 1<<GFBITS; i++ {
		a := bitrev(gf(i))

		for j := 0; j < GFBITS; j++ {
			L[j][i/8] |= byte(((a >> j) & 1) << (i % 8))
		}
	}

	for j := 0; j < GFBITS; j++ {
		apply_benes(L[j][:], c, 0)
	}

	for i := 0; i < SYS_N; i++ {
		s[i] = 0
		for j := GFBITS - 1; j >= 0; j-- {
			s[i] <<= 1
			s[i] |= gf((uint64(L[j][i/8]) >> (i % 8)) & 1)
		}
	}
}
