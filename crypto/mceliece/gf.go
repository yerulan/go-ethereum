// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package mceliece

type gf uint16

func gf_iszero(a gf) gf {
	var t uint32 = uint32(a)
	t -= 1
	t >>= 19
	return gf(t)
}

func gf_add(in0, in1 gf) gf {
	return in0 ^ in1
}

func gf_mul(in0, in1 gf) gf {
	var i int
	var tmp, t0, t1, t uint64

	t0 = uint64(in0)
	t1 = uint64(in1)

	tmp = t0 * (t1 & 1)

	for i = 1; i < GFBITS; i++ {
		tmp ^= (t0 * (t1 & (1 << i)))
	}

	t = tmp & 0x1FF0000
	tmp ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	t = tmp & 0x000E000
	tmp ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	return gf(tmp & GFMASK)
}

func gf_sq2(in gf) gf {
	var i int
	var B = [4]uint64{0x1111111111111111, 0x0303030303030303, 0x000F000F000F000F, 0x000000FF000000FF}
	var M = [4]uint64{0x0001FF0000000000, 0x000000FF80000000, 0x000000007FC00000, 0x00000000003FE000}

	var x uint64 = uint64(in)
	var t uint64

	x = (x | (x << 24)) & B[3]
	x = (x | (x << 12)) & B[2]
	x = (x | (x << 6)) & B[1]
	x = (x | (x << 3)) & B[0]

	for i = 0; i < 4; i++ {
		t = x & M[i]
		x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	}

	return gf(x & GFMASK)
}

func gf_sqmul(in, m gf) gf {
	var i int

	var x, t0, t1, t uint64

	var M = [3]uint64{0x0000001FF0000000, 0x000000000FF80000, 0x000000000007E000}

	t0 = uint64(in)
	t1 = uint64(m)

	x = (t1 << 6) * (t0 & (1 << 6))

	t0 ^= (t0 << 7)

	x ^= (t1 * (t0 & (0x04001)))
	x ^= (t1 * (t0 & (0x08002))) << 1
	x ^= (t1 * (t0 & (0x10004))) << 2
	x ^= (t1 * (t0 & (0x20008))) << 3
	x ^= (t1 * (t0 & (0x40010))) << 4
	x ^= (t1 * (t0 & (0x80020))) << 5
	for i = 0; i < 3; i++ {
		t = x & M[i]
		x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	}

	return gf(x & GFMASK)
}

func gf_sq2mul(in, m gf) gf {
	var i int

	var x, t0, t1, t uint64

	var M = [6]uint64{0x1FF0000000000000, 0x000FF80000000000, 0x000007FC00000000, 0x00000003FE000000, 0x0000000001FE0000, 0x000000000001E000}

	t0 = uint64(in)
	t1 = uint64(m)

	x = (t1 << 18) * (t0 & (1 << 6))

	t0 ^= (t0 << 21)

	x ^= (t1 * (t0 & (0x010000001)))
	x ^= (t1 * (t0 & (0x020000002))) << 3
	x ^= (t1 * (t0 & (0x040000004))) << 6
	x ^= (t1 * (t0 & (0x080000008))) << 9
	x ^= (t1 * (t0 & (0x100000010))) << 12
	x ^= (t1 * (t0 & (0x200000020))) << 15

	for i = 0; i < 6; i++ {
		t = x & M[i]
		x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	}

	return gf(x & GFMASK)
}

func gf_frac(den, num gf) gf {
	var tmp_11, tmp_1111, out gf

	tmp_11 = gf_sqmul(den, den)          // ^11
	tmp_1111 = gf_sq2mul(tmp_11, tmp_11) // ^1111
	out = gf_sq2(tmp_1111)
	out = gf_sq2mul(out, tmp_1111) // ^11111111
	out = gf_sq2(out)
	out = gf_sq2mul(out, tmp_1111) // ^111111111111

	return gf_sqmul(out, num) // ^1111111111110 = ^-1
}

func gf_inv(den gf) gf {
	return gf_frac(den, 1)
}

func GF_mul(out, in0, in1 []gf) {
	var i, j int

	prod := make([]gf, SYS_T*2-1)

	for i = 0; i < SYS_T*2-1; i++ {
		prod[i] = 0
	}

	for i = 0; i < SYS_T; i++ {
		for j = 0; j < SYS_T; j++ {
			prod[i+j] ^= gf_mul(in0[i], in1[j])
		}
	}

	for i = (SYS_T - 1) * 2; i >= SYS_T; i-- {
		prod[i-SYS_T+7] ^= prod[i]
		prod[i-SYS_T+2] ^= prod[i]
		prod[i-SYS_T+1] ^= prod[i]
		prod[i-SYS_T+0] ^= prod[i]
	}

	for i = 0; i < SYS_T; i++ {
		out[i] = prod[i]
	}
}
