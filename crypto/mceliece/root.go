package mceliece

func eval(f []gf, a gf) gf {
	var r gf
	r = f[SYS_T]

	for i := SYS_T - 1; i >= 0; i-- {
		r = gf_mul(r, a)
		r = gf_add(r, f[i])
	}

	return r
}

func root(out, f, L []gf) {
	for i := 0; i < SYS_N; i++ {
		out[i] = eval(f, L[i])
	}
}
