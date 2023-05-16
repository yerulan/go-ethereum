package mceliece

func decrypt(e, sk, c []byte) int {
	var w int
	var check uint16

	r := make([]byte, SYS_N/8)

	g := make([]gf, SYS_T+1)
	L := make([]gf, SYS_N)

	s := make([]gf, SYS_T*2)
	sCmp := make([]gf, SYS_T*2)
	locator := make([]gf, SYS_T+1)
	images := make([]gf, SYS_N)

	var t gf

	// Copy ciphertext to r
	copy(r[:SYND_BYTES], c[:SYND_BYTES])
	for i := SYND_BYTES; i < SYS_N/8; i++ {
		r[i] = 0
	}

	// Load g and set the last element to 1
	for i := 0; i < SYS_T; i++ {
		g[i] = gf(load_gf(sk))
		sk = sk[2:]
	}
	g[SYS_T] = 1

	support_gen(L, sk)

	synd(s, g, L, r)

	bm(locator, s)

	root(images, locator, L)

	for i := 0; i < SYS_N/8; i++ {
		e[i] = 0
	}

	for i := 0; i < SYS_N; i++ {
		t = gf_iszero(images[i]) & 1

		e[i/8] |= byte(t << (i % 8))
		w += int(t)
	}

	synd(sCmp, g, L, e)

	check = uint16(w) ^ SYS_T

	for i := 0; i < SYS_T*2; i++ {
		check |= uint16(s[i] ^ sCmp[i])
	}

	check -= 1
	check >>= 15

	return int(check ^ 1)
}
