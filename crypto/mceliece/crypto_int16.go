package mceliece

import (
	"unsafe"
)

type crypto_int16 int16

type crypto_uint16 uint16
type crypto_uint16_signed int16

func crypto_declassify(x unsafe.Pointer, n uintptr) {

}

func crypto_int16_negative_mask(crypto_int16_x crypto_int16) crypto_int16 {
	return crypto_int16_x >> 15
}

func crypto_int16_nonzero_mask(crypto_int16_x crypto_int16) crypto_int16 {
	return crypto_int16_negative_mask(crypto_int16_x) | crypto_int16_negative_mask(-crypto_int16_x)
}

func crypto_int16_zero_mask(crypto_int16_x crypto_int16) crypto_int16 {
	return ^crypto_int16_nonzero_mask(crypto_int16_x)
}

func crypto_int16_positive_mask(crypto_int16_x crypto_int16) crypto_int16 {
	crypto_int16_z := -crypto_int16_x
	crypto_int16_z ^= crypto_int16_x & crypto_int16_z
	return crypto_int16_negative_mask(crypto_int16_z)
}

func crypto_int16_unequal_mask(crypto_int16_x, crypto_int16_y crypto_int16) crypto_int16 {
	crypto_int16_xy := crypto_int16_x ^ crypto_int16_y
	return crypto_int16_nonzero_mask(crypto_int16_xy)
}

func crypto_int16_equal_mask(crypto_int16_x, crypto_int16_y crypto_int16) crypto_int16 {
	return ^crypto_int16_unequal_mask(crypto_int16_x, crypto_int16_y)
}

func crypto_int16_smaller_mask(crypto_int16_x, crypto_int16_y crypto_int16) crypto_int16 {
	crypto_int16_xy := crypto_int16_x ^ crypto_int16_y
	crypto_int16_z := crypto_int16_x - crypto_int16_y
	crypto_int16_z ^= crypto_int16_xy & (crypto_int16_z ^ crypto_int16_x)
	return crypto_int16_negative_mask(crypto_int16_z)
}

func crypto_int16_min(crypto_int16_x, crypto_int16_y crypto_int16) crypto_int16 {
	crypto_int16_xy := crypto_int16_y ^ crypto_int16_x
	crypto_int16_z := crypto_int16_y - crypto_int16_x
	crypto_int16_z ^= crypto_int16_xy & (crypto_int16_z ^ crypto_int16_y)
	crypto_int16_z = crypto_int16_negative_mask(crypto_int16_z)
	crypto_int16_z &= crypto_int16_xy
	return crypto_int16_x ^ crypto_int16_z
}

func crypto_int16_max(crypto_int16_x, crypto_int16_y crypto_int16) crypto_int16 {
	crypto_int16_xy := crypto_int16_y ^ crypto_int16_x
	crypto_int16_z := crypto_int16_y - crypto_int16_x
	crypto_int16_z ^= crypto_int16_xy & (crypto_int16_z ^ crypto_int16_y)
	crypto_int16_z = crypto_int16_negative_mask(crypto_int16_z)
	crypto_int16_z &= crypto_int16_xy
	return crypto_int16_y ^ crypto_int16_z
}

func crypto_int16_minmax(crypto_int16_a, crypto_int16_b *crypto_int16) {
	crypto_int16_x := *crypto_int16_a
	crypto_int16_y := *crypto_int16_b
	crypto_int16_xy := crypto_int16_y ^ crypto_int16_x
	crypto_int16_z := crypto_int16_y - crypto_int16_x
	crypto_int16_z ^= crypto_int16_xy & (crypto_int16_z ^ crypto_int16_y)
	crypto_int16_z = crypto_int16_negative_mask(crypto_int16_z)
	crypto_int16_z &= crypto_int16_xy
	*crypto_int16_a = crypto_int16_x ^ crypto_int16_z
	*crypto_int16_b = crypto_int16_y ^ crypto_int16_z
}

func crypto_uint16_signed_negative_mask(crypto_uint16_signed_x crypto_uint16_signed) crypto_uint16_signed {
	return crypto_uint16_signed_x >> 15
}

func crypto_uint16_nonzero_mask(crypto_uint16_x crypto_uint16) crypto_uint16 {
	return crypto_uint16(crypto_uint16_signed_negative_mask(crypto_uint16_signed(crypto_uint16_x)) | crypto_uint16_signed_negative_mask(-crypto_uint16_signed(crypto_uint16_x)))
}

func crypto_uint16_zero_mask(crypto_uint16_x crypto_uint16) crypto_uint16 {
	return ^crypto_uint16_nonzero_mask(crypto_uint16_x)
}

func crypto_uint16_unequal_mask(crypto_uint16_x crypto_uint16, crypto_uint16_y crypto_uint16) crypto_uint16 {
	crypto_uint16_xy := crypto_uint16_x ^ crypto_uint16_y
	return crypto_uint16_nonzero_mask(crypto_uint16_xy)
}

func crypto_uint16_equal_mask(crypto_uint16_x crypto_uint16, crypto_uint16_y crypto_uint16) crypto_uint16 {
	return ^crypto_uint16_unequal_mask(crypto_uint16_x, crypto_uint16_y)
}

func crypto_uint16_smaller_mask(crypto_uint16_x crypto_uint16, crypto_uint16_y crypto_uint16) crypto_uint16 {
	crypto_uint16_xy := crypto_uint16_x ^ crypto_uint16_y
	crypto_uint16_z := crypto_uint16_x - crypto_uint16_y
	crypto_uint16_z ^= crypto_uint16_xy & (crypto_uint16_z ^ crypto_uint16_x ^ (1 << 15))
	return crypto_uint16(crypto_uint16_signed_negative_mask(crypto_uint16_signed(crypto_uint16_z)))
}

func crypto_uint16_min(crypto_uint16_x crypto_uint16, crypto_uint16_y crypto_uint16) crypto_uint16 {
	crypto_uint16_xy := crypto_uint16_y ^ crypto_uint16_x
	crypto_uint16_z := crypto_uint16_y - crypto_uint16_x
	crypto_uint16_z ^= crypto_uint16_xy & (crypto_uint16_z ^ crypto_uint16_y ^ (1 << 15))
	crypto_uint16_z = crypto_uint16(crypto_uint16_signed_negative_mask(crypto_uint16_signed(crypto_uint16_z)))
	crypto_uint16_z &= crypto_uint16_xy
	return crypto_uint16_x ^ crypto_uint16_z
}

func crypto_uint16_max(crypto_uint16_x crypto_uint16, crypto_uint16_y crypto_uint16) crypto_uint16 {
	crypto_uint16_xy := crypto_uint16_y ^ crypto_uint16_x
	crypto_uint16_z := crypto_uint16_y - crypto_uint16_x
	crypto_uint16_z ^= crypto_uint16_xy & (crypto_uint16_z ^ crypto_uint16_y ^ (1 << 15))
	crypto_uint16_z = crypto_uint16(crypto_uint16_signed_negative_mask(crypto_uint16_signed(crypto_uint16_z)))
	crypto_uint16_z &= crypto_uint16_xy
	return crypto_uint16_y ^ crypto_uint16_z
}

func crypto_uint16_minmax(crypto_uint16_a *crypto_uint16, crypto_uint16_b *crypto_uint16) {
	crypto_uint16_x := *crypto_uint16_a
	crypto_uint16_y := *crypto_uint16_b
	crypto_uint16_xy := crypto_uint16_y ^ crypto_uint16_x
	crypto_uint16_z := crypto_uint16_y - crypto_uint16_x
	crypto_uint16_z ^= crypto_uint16_xy & (crypto_uint16_z ^ crypto_uint16_y ^ (1 << 15))
	crypto_uint16_z = crypto_uint16(crypto_uint16_signed_negative_mask(crypto_uint16_signed(crypto_uint16_z)))
	crypto_uint16_z &= crypto_uint16_xy
	*crypto_uint16_a = crypto_uint16_x ^ crypto_uint16_z
	*crypto_uint16_b = crypto_uint16_y ^ crypto_uint16_z
}
