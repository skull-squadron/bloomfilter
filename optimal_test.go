package bloomfilter

import (
	"testing"
)

func TestOptimal(t *testing.T) {
	tests := []struct {
		n    uint64
		p    float64
		k, m uint64
	}{
		{
			n: 1000,
			p: 0.01 / 100,
			k: 14,
			m: 19171,
		},
		{
			n: 10000,
			p: 0.01 / 100,
			k: 14,
			m: 191702,
		},
		{
			n: 10000,
			p: 0.01 / 100,
			k: 14,
			m: 191702,
		},
		{
			n: 1000,
			p: 0.001 / 100,
			k: 17,
			m: 23963,
		},
	}

	for _, test := range tests {
		m := OptimalM(test.n, test.p)
		k := OptimalK(m, test.n)

		if k != test.k || m != test.m {
			t.Errorf(
				"n=%d p=%f: expected (m=%d, k=%d), got (m=%d, k=%d)",
				test.n,
				test.p,
				test.m,
				test.k,
				m,
				k,
			)
		}
	}
}
