package mceliece

import "unsafe"

func gf_is_zero_declassify(t gf) uint16 {
	mask := crypto_uint16_zero_mask(crypto_uint16(t))
	crypto_declassify(unsafe.Pointer(&mask), unsafe.Sizeof(mask))
	return uint16(mask)
}

func genpoly_gen(out, f []gf) int {
	var i, j, k, c int
	mat := make([][]gf, SYS_T+1)
	var mask, inv, t gf

	// fill matrix
	mat[0][0] = 1

	for i = 1; i < SYS_T; i++ {
		mat[0][i] = 0
	}

	for i = 0; i < SYS_T; i++ {
		mat[1][i] = f[i]
	}

	for j = 2; j <= SYS_T; j++ {
		GF_mul(mat[j], mat[j-1], f)
	}

	// gaussian
	for j = 0; j < SYS_T; j++ {
		for k = j + 1; k < SYS_T; k++ {
			mask = gf_iszero(mat[j][j])

			for c = j; c < SYS_T+1; c++ {
				mat[c][j] ^= mat[c][k] & mask
			}
		}

		if gf_is_zero_declassify(mat[j][j]) != 0 { // return if not systematic
			return -1
		}

		inv = gf_inv(mat[j][j])

		for c = j; c < SYS_T+1; c++ {
			mat[c][j] = gf_mul(mat[c][j], inv)
		}

		for k = 0; k < SYS_T; k++ {
			if k != j {
				t = mat[j][k]

				for c = j; c < SYS_T+1; c++ {
					mat[c][k] ^= gf_mul(mat[c][j], t)
				}
			}
		}
	}

	for i = 0; i < SYS_T; i++ {
		out[i] = mat[SYS_T][i]
	}

	return 0
}
