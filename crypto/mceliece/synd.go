package mceliece

func synd(out, f, L []gf, r []byte) {
	for j := 0; j < 2*SYS_T; j++ {
		out[j] = 0
	}

	for i := 0; i < SYS_N; i++ {
		c := (r[i/8] >> (i % 8)) & 1

		e := eval(f, L[i])
		eInv := gf_inv(gf_mul(e, e))

		for j := 0; j < 2*SYS_T; j++ {
			out[j] = gf_add(out[j], gf_mul(eInv, gf(c)))
			eInv = gf_mul(eInv, L[i])
		}
	}
}
