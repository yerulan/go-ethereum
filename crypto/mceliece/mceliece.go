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

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto/blake2b"
)

// Define the McEliece cryptosystem parameters
const (
	SYS_N                  = 6688
	SYS_T                  = 128
	GFBITS                 = 13
	GFMASK                 = (1 << GFBITS) - 1
	CRYPTO_PUBLICKEYBYTES  = 1044992
	CRYPTO_SECRETKEYBYTES  = 13932
	CRYPTO_CIPHERTEXTBYTES = 208
	CRYPTO_BYTES           = 32
	COND_BYTES             = (1 << (GFBITS - 4)) * (2*GFBITS - 1)
	IRR_BYTES              = SYS_T * 2
	PK_NROWS               = SYS_T * GFBITS
	PK_NCOLS               = SYS_N - PK_NROWS
	PK_ROW_BYTES           = (PK_NCOLS + 7) / 8
	SYND_BYTES             = (PK_NROWS + 7) / 8
)

// PublicKey represents a McEliece public key
type PublicKey struct {
	G []byte // Generator matrix
}

// PrivateKey represents a McEliece private key
type PrivateKey struct {
	S []byte // Scrambling matrix
	P []byte // Permutation matrix
}

func Encrypt(c, key, pk []byte) int {
	e := make([]byte, SYS_N/8)
	oneEC := make([]byte, 1+SYS_N/8+SYND_BYTES)
	oneEC[0] = 1

	encrypt(c, pk, e)

	copy(oneEC[1:], e)
	copy(oneEC[1+SYS_N/8:], c[:SYND_BYTES])

	key = blake2b.Sum256(oneEC)

	return 0
}

func Decrypt(key, c, sk []byte) int {
	var retDecrypt byte

	var m uint16

	e := make([]byte, SYS_N/8)
	preimage := make([]byte, 1+SYS_N/8+SYND_BYTES)
	x := preimage
	s := sk[40+IRR_BYTES+COND_BYTES:]

	retDecrypt = byte(decrypt(e, sk[40:], c))

	m = uint16(retDecrypt)
	m -= 1
	m >>= 8

	x[0] = byte(m & 1)
	x = x[1:]
	for i := 0; i < SYS_N/8; i++ {
		x[0] = (^byte(m) & s[i]) | (byte(m) & e[i])
		x = x[1:]
	}
	for i := 0; i < SYND_BYTES; i++ {
		x[0] = c[i]
		x = x[1:]
	}

	key = blake2b.Sum256(preimage)

	return 0
}

func GenerateKeys(pk, sk []byte) int {
	seed := make([]byte, 33)
	r := make([]byte, SYS_N/8+(1<<GFBITS)*4+SYS_T*2+32)
	f := make([]gf, SYS_T)
	irr := make([]gf, SYS_T)
	perm := make([]uint32, 1<<GFBITS)
	pi := make([]int16, 1<<GFBITS)
	seed[0] = 64

	skp := sk
	rp := r[len(r)-32:]

	for {
		rp = r[len(r)-32:]
		skp = sk

		r = blake2b.Sum256(seed)
		copy(skp, seed[1:33])
		skp = skp[32+8:]
		copy(seed[1:], rp)

		rp = rp[:len(f)]
		for i := 0; i < SYS_T; i++ {
			f[i] = gf(load_gf(rp[i*2:]))
		}

		if genpoly_gen(irr, f) != 0 {
			continue
		}

		for i := 0; i < SYS_T; i++ {
			store_gf(skp[i*2:], uint16(irr[i]))
		}

		rp = rp[:len(perm)*4]
		for i := 0; i < len(perm); i++ {
			perm[i] = binary.LittleEndian.Uint32(rp[i*4:])
		}

		if pk_gen(pk, skp, perm, pi) != 0 {
			continue
		}

		skp = skp[IRR_BYTES:]

		controlbitsfrompermutation(skp, pi, GFBITS, 1<<GFBITS)
		skp = skp[COND_BYTES:]

		rp = rp[:SYS_N/8]
		copy(skp, rp)

		binary.LittleEndian.PutUint64(sk[32:], 0xFFFFFFFFFFFFFFFF)

		break
	}

	return 0
}
