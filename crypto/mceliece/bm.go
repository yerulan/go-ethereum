package mceliece

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bm(out, s []gf) {
	N := uint16(0)
	L := uint16(0)
	var mle, mne uint16

	T := make([]gf, SYS_T+1)
	C := make([]gf, SYS_T+1)
	B := make([]gf, SYS_T+1)

	b := gf(1)
	var d, f gf

	for i := 0; i < SYS_T+1; i++ {
		C[i] = 0
		B[i] = 0
	}

	B[1] = 1
	C[0] = 1

	for N = 0; N < 2*SYS_T; N++ {
		d = 0

		for i := 0; i <= min(int(N), SYS_T); i++ {
			d ^= gf_mul(C[i], s[N-uint16(i)])
		}

		mne = uint16(d)
		mne -= 1
		mne >>= 15
		mne -= 1

		mle = N
		mle -= 2 * L
		mle >>= 15
		mle -= 1

		mle &= mne

		for i := 0; i <= SYS_T; i++ {
			T[i] = C[i]
		}

		f = gf_frac(b, d)

		for i := 0; i <= SYS_T; i++ {
			C[i] ^= gf_mul(f, B[i]) & gf(mne)
		}

		L = (L & ^mle) | ((N + 1 - L) & mle)

		for i := 0; i <= SYS_T; i++ {
			B[i] = (B[i] & ^gf(mle)) | (T[i] & gf(mle))
		}

		b = (b & ^gf(mle)) | (d & gf(mle))

		for i := SYS_T; i >= 1; i-- {
			B[i] = B[i-1]
		}
		B[0] = 0
	}

	for i := 0; i <= SYS_T; i++ {
		out[i] = C[SYS_T-i]
	}
}
