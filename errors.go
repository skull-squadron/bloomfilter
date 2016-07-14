//
// Face-meltingly fast, thread-safe, marshalable, unionable, probability- and optimal-size-calculating Bloom filter in go
//
// https://github.com/steakknife/bloomfilter
//
// Copyright © 2014, 2015 Barry Allard
//
// MIT license
//
package bloomfilter

import (
	"fmt"
)

var (
	hashError                     = fmt.Errorf("Hash mismatch, the Bloom filter is probably corrupt.")
	kError                        = fmt.Errorf("keys must have length %d or greater", K_MIN)
	mError                        = fmt.Errorf("m (number of bits in the Bloom filter) must be >= %d", M_MIN)
	uniqueKeysError               = fmt.Errorf("Bloom filter keys must be unique")
	incompatibleBloomFiltersError = fmt.Errorf("Cannot perform union on two incompatible Bloom filters")
)
