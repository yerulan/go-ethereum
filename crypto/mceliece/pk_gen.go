package mceliece

import "unsafe"

func uint64_is_equal_declassify(t uint64, u uint64) uint64 {
	mask := crypto_uint64_equal_mask(crypto_uint64(t), crypto_uint64(u))
	crypto_declassify(unsafe.Pointer(&mask), unsafe.Sizeof(mask))
	return uint64(mask)
}

func uint64_is_zero_declassify(t uint64) uint64 {
	mask := crypto_uint64_zero_mask(crypto_uint64(t))
	crypto_declassify(unsafe.Pointer(&mask), unsafe.Sizeof(mask))
	return uint64(mask)
}

func pk_gen(pk []byte, sk []byte, perm []uint32, pi []int16) int {
	var row int
	var mask byte

	g := make([]gf, SYS_T+1) // Goppa polynomial
	L := make([]gf, SYS_N)   // support
	inv := make([]gf, SYS_N)

	g[SYS_T] = 1

	for i := 0; i < SYS_T; i++ {
		g[i] = gf(load_gf(sk))
		sk = sk[2:]
	}

	buf := make([]uint64, 1<<GFBITS)

	for i := 0; i < (1 << GFBITS); i++ {
		buf[i] = uint64(perm[i])
		buf[i] <<= 31
		buf[i] |= uint64(i)
	}

	uint64_sort(buf, 1<<GFBITS)

	for i := 1; i < (1 << GFBITS); i++ {
		if uint64_is_equal_declassify(buf[i-1]>>31, buf[i]>>31) != 0 {
			return -1
		}
	}

	for i := 0; i < (1 << GFBITS); i++ {
		pi[i] = int16(buf[i] & GFMASK)
	}
	for i := 0; i < SYS_N; i++ {
		L[i] = bitrev(gf(pi[i]))
	}

	// filling the matrix

	root(inv, g, L)

	for i := 0; i < SYS_N; i++ {
		inv[i] = gf_inv(inv[i])
	}

	mat := make([][]byte, PK_NROWS)
	for i := 0; i < PK_NROWS; i++ {
		mat[i] = make([]byte, SYS_N/8)
	}

	for i := 0; i < SYS_T; i++ {
		for j, k := 0, 0; j < SYS_N; j, k = j+8, k+1 {
			for b := 0; b < GFBITS; b++ {
				bit := ((inv[j+7] >> b) & 1) << 7
				bit |= ((inv[j+6] >> b) & 1) << 6
				bit |= ((inv[j+5] >> b) & 1) << 5
				bit |= ((inv[j+4] >> b) & 1) << 4
				bit |= ((inv[j+3] >> b) & 1) << 3
				bit |= ((inv[j+2] >> b) & 1) << 2
				bit |= ((inv[j+1] >> b) & 1) << 1
				bit |= (inv[j+0] >> b) & 1
				mat[i*GFBITS+b][k] = byte(bit)
			}

			for j := 0; j < SYS_N; j++ {
				inv[j] = gf_mul(inv[j], L[j])
			}
		}
	}

	// gaussian elimination

	for i := 0; i < (PK_NROWS+7)/8; i++ {
		for j := 0; j < 8; j++ {
			row = i*8 + j

			if row >= PK_NROWS {
				break
			}

			for k := row + 1; k < PK_NROWS; k++ {
				mask = mat[row][i] ^ mat[k][i]
				mask >>= j
				mask &= 1
				mask = -mask

				for c := 0; c < SYS_N/8; c++ {
					mat[row][c] ^= mat[k][c] & mask
				}
			}

			if uint64_is_zero_declassify((uint64(mat[row][i])>>j)&1) != 0 { // return if not systematic
				return -1
			}

			for k := 0; k < PK_NROWS; k++ {
				if k != row {
					mask = mat[k][i] >> j
					mask &= 1
					mask = -mask

					for c := 0; c < SYS_N/8; c++ {
						mat[k][c] ^= mat[row][c] & mask
					}
				}
			}
		}
	}

	for i := 0; i < PK_NROWS; i++ {
		copy(pk[i*PK_ROW_BYTES:], mat[i][PK_NROWS/8:])
	}

	return 0
}
