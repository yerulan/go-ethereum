package mceliece

import "unsafe"

type int16_t int16
type int32_t int32

func int32_min(a, b int32_t) int32_t {
	if a < b {
		return a
	}
	return b
}

func int32_sort(arr []int32_t, n int) {
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j = j - 1
		}
		arr[j+1] = key
	}
}

func cbrecursion(out []byte, pos, step int64, pi []int16, w, n int64, temp []int32_t) {
	A := temp[:n]
	B := temp[n : n+n]
	q := (*[1 << 20]int16)(unsafe.Pointer(&temp[n+n/4]))[:n]

	if w == 1 {
		out[pos>>3] ^= byte(pi[0] << (pos & 7))
		return
	}

	for x := int64(0); x < n; x++ {
		A[x] = int32_t(((pi[x] ^ 1) << 16) | pi[x^1])
	}
	int32_sort(A, int(n))

	for x := int64(0); x < n; x++ {
		Ax := A[x]
		px := Ax & 0xffff
		cx := int32_min(px, int32_t(x))
		B[x] = ((px << 16) | cx)
	}

	for x := int64(0); x < n; x++ {
		A[x] = (A[x] << 16) | int32_t(x)
	}
	int32_sort(A, int(n))

	for x := int64(0); x < n; x++ {
		A[x] = (A[x] << 16) + (B[x] >> 16)
	}
	int32_sort(A, int(n))

	if w <= 10 {
		for x := int64(0); x < n; x++ {
			B[x] = (((A[x] & 0xffff) << 10) | (B[x] & 0x3ff))
		}

		for i := int64(1); i < w-1; i++ {
			for x := int64(0); x < n; x++ {
				A[x] = (((B[x] &^ 0x3ff) << 6) | int32_t(x))
			}
			int32_sort(A, int(n))

			for x := int64(0); x < n; x++ {
				A[x] = (A[x] << 20) | B[x]
			}
			int32_sort(A, int(n))

			for x := int64(0); x < n; x++ {
				ppcpx := A[x] & 0xfffff
				ppcx := ((A[x] & 0xffc00) | (B[x] & 0x3ff))
				B[x] = int32_min(ppcx, ppcpx)
			}
		}

		for x := int64(0); x < n; x++ {
			B[x] &= 0x3ff
		}
	} else {
		for x := int64(0); x < n; x++ {
			B[x] = ((A[x] << 16) | (B[x] & 0xffff))
		}

		for i := int64(1); i < w-1; i++ {
			for x := int64(0); x < n; x++ {
				A[x] = (B[x] & ^0xffff) | int32_t(x)
			}
			int32_sort(A, int(n))
			for x := int64(0); x < n; x++ {
				A[x] = (A[x] << 16) | (B[x] & 0xffff)
			}

			if i < w-2 {
				for x := int64(0); x < n; x++ {
					B[x] = (A[x] &^ 0xffff) | (B[x] >> 16)
				}
				int32_sort(B, int(n))

				for x := int64(0); x < n; x++ {
					B[x] = (B[x] << 16) | (A[x] & 0xffff)
				}
			}

			int32_sort(A, int(n))

			for x := int64(0); x < n; x++ {
				cpx := (B[x] &^ 0xffff) | (A[x] & 0xffff)
				B[x] = int32_min(B[x], cpx)
			}
		}

		for x := int64(0); x < n; x++ {
			B[x] &= 0xffff
		}
	}
	for x := int64(0); x < n; x++ {
		A[x] = (int32_t(pi[x]) << 16) + int32_t(x)
	}
	int32_sort(A, int(n))
	pos += (2*w - 3) * step * (n / 2)

	for k := int64(0); k < n/2; k++ {
		y := 2 * k
		lk := B[y] & 1        /* l[k] */
		Ly := int32_t(y) + lk /* L[y] */
		Ly1 := Ly ^ 1         /* L[y+1] */

		out[pos>>3] ^= byte(lk << (pos & 7))
		pos += step

		A[y] = (Ly << 16) | (B[y] & 0xffff)
		A[y+1] = (Ly1 << 16) | (B[y+1] & 0xffff)
	}

	int32_sort(A, int(n)) /* A = (id<<16)+F(pi(L)) = (id<<16)+M */

	pos -= (2*w - 2) * step * (n / 2)

	for j := int64(0); j < n/2; j++ {
		q[j] = int16(A[2*j]&0xffff) >> 1
		q[j+n/2] = int16(A[2*j+1]&0xffff) >> 1
	}

	cbrecursion(out, pos, step*2, q, w-1, n/2, temp)
	cbrecursion(out, pos+step, step*2, q[n/2:], w-1, n/2, temp)
}

func layer(p []int16_t, cb []byte, s, n int) {
	stride := 1 << s
	index := 0
	var d, m int16_t

	for i := 0; i < n; i += stride * 2 {
		for j := 0; j < stride; j++ {
			d = p[i+j] ^ p[i+j+stride]
			m = int16_t(((cb[index>>3]) >> (index & 7)) & 1)
			m = -m
			d &= m
			p[i+j] ^= d
			p[i+j+stride] ^= d
			index++
		}
	}
}

func controlbitsfrompermutation(out []byte, pi []int16, w, n int64) {
	temp := make([]int32_t, 2*n)
	for {
		copy(out, make([]byte, len(out)))
		cbrecursion(out, 0, 1, pi, w, n, temp)

		// check for correctness
		piTest := make([]int16_t, n)
		for i := int64(0); i < n; i++ {
			piTest[i] = int16_t(i)
		}

		ptr := out
		for i := int64(0); i < w; i++ {
			layer(piTest, ptr, int(i), int(n))
			ptr = ptr[n>>4:]
		}

		for i := w - 2; i >= 0; i-- {
			layer(piTest, ptr, int(i), int(n))
			ptr = ptr[n>>4:]
		}

		diff := int16_t(0)
		for i := int64(0); i < n; i++ {
			diff |= int16_t(pi[i]) ^ piTest[i]
		}

		diff = int16_t(crypto_int16_nonzero_mask(crypto_int16(diff)))
		if diff == 0 {
			break
		}
	}
}
