package mceliece

import (
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

const (
	RNG_SUCCESS     = 0
	RNG_BAD_MAXLEN  = -1
	RNG_BAD_OUTBUF  = -2
	RNG_BAD_REQ_LEN = -3
)

var DRBG_ctx AES256_CTR_DRBG_struct

type AES_XOF_struct struct {
	buffer           [16]byte
	buffer_pos       int
	length_remaining uint64
	key              [32]byte
	ctr              [16]byte
}

type AES256_CTR_DRBG_struct struct {
	Key            [32]byte
	V              [16]byte
	reseed_counter int
}

func AES256_ECB(key, ctr, buffer []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating AES cipher:", err)
		os.Exit(1)
	}

	if len(ctr) != block.BlockSize() || len(buffer) != block.BlockSize() {
		fmt.Fprintln(os.Stderr, "Key or input size mismatch")
		os.Exit(1)
	}

	block.Encrypt(buffer, ctr)
}

func seedexpander_init(ctx *AES_XOF_struct, seed, diversifier []byte, maxlen uint64) error {
	if maxlen >= 0x100000000 {
		return errors.New("RNG_BAD_MAXLEN")
	}

	ctx.length_remaining = maxlen

	copy(ctx.key[:], seed[:32])

	copy(ctx.ctr[:8], diversifier)
	binary.LittleEndian.PutUint32(ctx.ctr[8:12], uint32(maxlen))

	ctx.buffer_pos = 16

	return nil
}

func seedexpander(ctx *AES_XOF_struct, x []byte, xlen uint64) error {
	offset := uint64(0)
	for xlen > 0 {
		if x == nil {
			return errors.New("RNG_BAD_OUTBUF")
		}
		if xlen >= ctx.length_remaining {
			return errors.New("RNG_BAD_REQ_LEN")
		}

		ctx.length_remaining -= xlen

		for xlen > 0 {
			if xlen <= uint64(16-ctx.buffer_pos) {
				copy(x[offset:], ctx.buffer[ctx.buffer_pos:uint64(ctx.buffer_pos)+xlen])
				ctx.buffer_pos += int(xlen)

				return nil
			}

			copy(x[offset:], ctx.buffer[ctx.buffer_pos:])
			xlen -= uint64(16 - ctx.buffer_pos)
			offset += uint64(16 - ctx.buffer_pos)

			AES256_ECB(ctx.key[:], ctx.ctr[:], ctx.buffer[:])
			ctx.buffer_pos = 0

			for i := 15; i >= 12; i-- {
				if ctx.ctr[i] == 0xff {
					ctx.ctr[i] = 0x00
				} else {
					ctx.ctr[i]++
					break
				}
			}
		}
	}

	return nil
}

func handleErrors() {
	fmt.Fprintln(os.Stderr, "Error occurred")
	os.Exit(1)
}

func randombytes_init(entropy_input, personalization_string []byte, security_strength int) {
	var seed_material [48]byte

	copy(seed_material[:], entropy_input[:48])
	if personalization_string != nil {
		for i := 0; i < 48; i++ {
			seed_material[i] ^= personalization_string[i]
		}
	}

	DRBG_ctx := AES256_CTR_DRBG_struct{}
	copy(DRBG_ctx.Key[:], make([]byte, 32))
	copy(DRBG_ctx.V[:], make([]byte, 16))
	AES256_CTR_DRBG_Update(seed_material[:], DRBG_ctx.Key[:], DRBG_ctx.V[:])
	DRBG_ctx.reseed_counter = 1
}

func randombytes(x []byte, xlen uint64) error {
	i := 0
	for xlen > 0 {
		block := [16]byte{}
		for j := 15; j >= 0; j-- {
			if DRBG_ctx.V[j] == 0xff {
				DRBG_ctx.V[j] = 0x00
			} else {
				DRBG_ctx.V[j]++
				break
			}
		}
		AES256_ECB(DRBG_ctx.Key[:], DRBG_ctx.V[:], block[:])
		if x == nil {
			return errors.New("RNG_BAD_OUTBUF")
		}
		if xlen > 15 {
			copy(x[i:], block[:])
			i += 16
			xlen -= 16
		} else {
			copy(x[i:], block[:xlen])
			xlen = 0
		}
	}
	AES256_CTR_DRBG_Update(nil, DRBG_ctx.Key[:], DRBG_ctx.V[:])
	DRBG_ctx.reseed_counter++
	return nil
}

func AES256_CTR_DRBG_Update(provided_data, Key, V []byte) {
	temp := [48]byte{}
	for i := 0; i < 3; i++ {
		for j := 15; j >= 0; j-- {
			if V[j] == 0xff {
				V[j] = 0x00
			} else {
				V[j]++
				break
			}
		}

		AES256_ECB(Key, V, temp[i*16:(i+1)*16])
	}

	if provided_data != nil {
		for i := 0; i < 48; i++ {
			temp[i] ^= provided_data[i]
		}
	}

	copy(Key[:], temp[:32])
	copy(V[:], temp[32:])
}
