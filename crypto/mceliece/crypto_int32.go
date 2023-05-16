package mceliece

type crypto_int32 int32

type crypto_uint32 uint32
type crypto_uint32_signed int32

func crypto_uint32_signed_negative_mask(crypto_uint32_signed_x crypto_uint32_signed) crypto_uint32_signed {
	return crypto_uint32_signed_x >> 31
}

func crypto_uint32_nonzero_mask(crypto_uint32_x crypto_uint32) crypto_uint32 {
	return crypto_uint32(crypto_uint32_signed_negative_mask(crypto_uint32_signed(crypto_uint32_x)) | crypto_uint32_signed_negative_mask(-crypto_uint32_signed(crypto_uint32_x)))
}

func crypto_uint32_zero_mask(crypto_uint32_x crypto_uint32) crypto_uint32 {
	return ^crypto_uint32_nonzero_mask(crypto_uint32_x)
}

func crypto_uint32_unequal_mask(crypto_uint32_x, crypto_uint32_y crypto_uint32) crypto_uint32 {
	crypto_uint32_xy := crypto_uint32_x ^ crypto_uint32_y
	return crypto_uint32_nonzero_mask(crypto_uint32_xy)
}

func crypto_uint32_equal_mask(crypto_uint32_x, crypto_uint32_y crypto_uint32) crypto_uint32 {
	return ^crypto_uint32_unequal_mask(crypto_uint32_x, crypto_uint32_y)
}

func crypto_uint32_smaller_mask(crypto_uint32_x, crypto_uint32_y crypto_uint32) crypto_uint32 {
	crypto_uint32_xy := crypto_uint32_x ^ crypto_uint32_y
	crypto_uint32_z := crypto_uint32_x - crypto_uint32_y
	crypto_uint32_z ^= crypto_uint32_xy & (crypto_uint32_z ^ crypto_uint32_x ^ (crypto_uint32(1) << 31))
	return crypto_uint32(crypto_uint32_signed_negative_mask(crypto_uint32_signed(crypto_uint32_z)))
}

func crypto_uint32_min(crypto_uint32_x, crypto_uint32_y crypto_uint32) crypto_uint32 {
	crypto_uint32_xy := crypto_uint32_y ^ crypto_uint32_x
	crypto_uint32_z := crypto_uint32_y - crypto_uint32_x
	crypto_uint32_z ^= crypto_uint32_xy & (crypto_uint32_z ^ crypto_uint32_y ^ (crypto_uint32(1) << 31))
	crypto_uint32_z = crypto_uint32(crypto_uint32_signed_negative_mask(crypto_uint32_signed(crypto_uint32_z)))
	crypto_uint32_z &= crypto_uint32_xy
	return crypto_uint32_x ^ crypto_uint32_z
}

func crypto_uint32_max(crypto_uint32_x, crypto_uint32_y crypto_uint32) crypto_uint32 {
	crypto_uint32_xy := crypto_uint32_y ^ crypto_uint32_x
	crypto_uint32_z := crypto_uint32_y - crypto_uint32_x
	crypto_uint32_z ^= crypto_uint32_xy & (crypto_uint32_z ^ crypto_uint32_y ^ (crypto_uint32(1) << 31))
	crypto_uint32_z = crypto_uint32(crypto_uint32_signed_negative_mask(crypto_uint32_signed(crypto_uint32_z)))
	crypto_uint32_z &= crypto_uint32_xy
	return crypto_uint32_y ^ crypto_uint32_z
}

func crypto_uint32_minmax(crypto_uint32_a, crypto_uint32_b *crypto_uint32) {
	crypto_uint32_x := *crypto_uint32_a
	crypto_uint32_y := *crypto_uint32_b
	crypto_uint32_xy := crypto_uint32_y ^ crypto_uint32_x
	crypto_uint32_z := crypto_uint32_y - crypto_uint32_x
	crypto_uint32_z ^= crypto_uint32_xy & (crypto_uint32_z ^ crypto_uint32_y ^ (crypto_uint32(1) << 31))
	crypto_uint32_z = crypto_uint32(crypto_uint32_signed_negative_mask(crypto_uint32_signed(crypto_uint32_z)))
	crypto_uint32_z &= crypto_uint32_xy
	*crypto_uint32_a = crypto_uint32_x ^ crypto_uint32_z
	*crypto_uint32_b = crypto_uint32_y ^ crypto_uint32_z
}

func crypto_int32_negative_mask(crypto_int32_x crypto_int32) crypto_int32 {
	return crypto_int32_x >> 31
}

func crypto_int32_nonzero_mask(crypto_int32_x crypto_int32) crypto_int32 {
	return crypto_int32_negative_mask(crypto_int32_x) | crypto_int32_negative_mask(-crypto_int32_x)
}

func crypto_int32_zero_mask(crypto_int32_x crypto_int32) crypto_int32 {
	return ^crypto_int32_nonzero_mask(crypto_int32_x)
}

func crypto_int32_positive_mask(crypto_int32_x crypto_int32) crypto_int32 {
	crypto_int32_z := -crypto_int32_x
	crypto_int32_z ^= crypto_int32_x & crypto_int32_z
	return crypto_int32_negative_mask(crypto_int32_z)
}

func crypto_int32_unequal_mask(crypto_int32_x crypto_int32, crypto_int32_y crypto_int32) crypto_int32 {
	crypto_int32_xy := crypto_int32_x ^ crypto_int32_y
	return crypto_int32_nonzero_mask(crypto_int32_xy)
}

func crypto_int32_equal_mask(crypto_int32_x crypto_int32, crypto_int32_y crypto_int32) crypto_int32 {
	return ^crypto_int32_unequal_mask(crypto_int32_x, crypto_int32_y)
}

func crypto_int32_smaller_mask(crypto_int32_x crypto_int32, crypto_int32_y crypto_int32) crypto_int32 {
	crypto_int32_xy := crypto_int32_x ^ crypto_int32_y
	crypto_int32_z := crypto_int32_x - crypto_int32_y
	crypto_int32_z ^= crypto_int32_xy & (crypto_int32_z ^ crypto_int32_x)
	return crypto_int32_negative_mask(crypto_int32_z)
}

func crypto_int32_min(crypto_int32_x crypto_int32, crypto_int32_y crypto_int32) crypto_int32 {
	crypto_int32_xy := crypto_int32_y ^ crypto_int32_x
	crypto_int32_z := crypto_int32_y - crypto_int32_x
	crypto_int32_z ^= crypto_int32_xy & (crypto_int32_z ^ crypto_int32_y)
	crypto_int32_z = crypto_int32_negative_mask(crypto_int32_z)
	crypto_int32_z &= crypto_int32_xy
	return crypto_int32_x ^ crypto_int32_z
}

func crypto_int32_max(crypto_int32_x crypto_int32, crypto_int32_y crypto_int32) crypto_int32 {
	crypto_int32_xy := crypto_int32_y ^ crypto_int32_x
	crypto_int32_z := crypto_int32_y - crypto_int32_x
	crypto_int32_z ^= crypto_int32_xy & (crypto_int32_z ^ crypto_int32_y)
	crypto_int32_z = crypto_int32_negative_mask(crypto_int32_z)
	crypto_int32_z &= crypto_int32_xy
	return crypto_int32_y ^ crypto_int32_z
}

func crypto_int32_minmax(crypto_int32_a *crypto_int32, crypto_int32_b *crypto_int32) {
	crypto_int32_x := *crypto_int32_a
	crypto_int32_y := *crypto_int32_b
	crypto_int32_xy := crypto_int32_y ^ crypto_int32_x
	crypto_int32_z := crypto_int32_y - crypto_int32_x
	crypto_int32_z ^= crypto_int32_xy & (crypto_int32_z ^ crypto_int32_y)
	crypto_int32_z = crypto_int32_negative_mask(crypto_int32_z)
	crypto_int32_z &= crypto_int32_xy
	*crypto_int32_a = crypto_int32_x ^ crypto_int32_z
	*crypto_int32_b = crypto_int32_y ^ crypto_int32_z
}
