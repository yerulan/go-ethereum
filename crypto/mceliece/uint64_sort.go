package mceliece

func uint64_sort(x []uint64, n int64) {
	if n < 2 {
		return
	}
	top := int64(1)
	for top < n-top {
		top += top
	}

	for p := top; p > 0; p >>= 1 {
		for i := int64(0); i < n-p; i++ {
			if i&p == 0 {
				uint64MinMax(&x[i], &x[i+p])
			}
		}
		i := int64(0)
		for q := top; q > p; q >>= 1 {
			for ; i < n-q; i++ {
				if i&p == 0 {
					a := x[i+p]
					for r := q; r > p; r >>= 1 {
						uint64MinMax(&a, &x[i+r])
					}
					x[i+p] = a
				}
			}
		}
	}
}

func uint64MinMax(a, b *uint64) {
	c := *b - *a
	c >>= 63
	c = -c
	c &= *a ^ *b
	*a ^= c
	*b ^= c
}
