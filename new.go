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
	"crypto/rand"
	"encoding/binary"
)

const (
	M_MIN        = 2 // minimum Bloom filter bits count
	K_MIN        = 1 // minimum number of keys
	UINT64_BYTES = 8
)

// new with CSPRNG keys
//
// m is the size of the Bloom filter, in bits, >= 2
//
// k is the number of keys, >= 1
func New(m, k uint64) *Filter {
	return NewWithKeys(m, newRandKeys(k))
}

func newRandKeys(k uint64) []uint64 {
	keys := make([]uint64, k)
	err := binary.Read(rand.Reader, binary.LittleEndian, keys)
	if err != nil {
		panicf("Cannot read %d bytes from CSRPNG crypto/rand.Read (err=%v)", UINT64_BYTES, err)
	}
	return keys
}

// Create a new Filter compatible with f
func (f Filter) NewCompatible() *Filter {
	return NewWithKeys(f.m, f.keys)
}

// New optimal Bloom filter with CSPRNG keys
func NewOptimal(maxN uint64, p float64) *Filter {
	m := OptimalM(maxN, p)
	k := OptimalK(m, maxN)
	debug("New optimal bloom filter :: requested max elements (n):%d, probability of collision (p):%1.10f -> recommends -> bits (m): %d (%f GiB), number of keys (k): %d", maxN, p, m, float64(m)/(gigabitsPerGiB), k)
	return New(m, k)
}

func UniqueKeys(keys []uint64) bool {
	for j := 0; j < len(keys)-1; j++ {
		elem := keys[j]
		for i := 1; i < j; i++ {
			if keys[i] == elem {
				return false
			}
		}
	}
	return true
}

// new with user-supplied keys
func NewWithKeys(m uint64, keys []uint64) *Filter {
	if m < M_MIN {
		panic(errM)
	}
	if len(keys) < K_MIN {
		panic(errK)
	}
	if !UniqueKeys(keys) {
		panic(errUniqueKeys)
	}
	return &Filter{
		m:    m,
		n:    0,
		bits: make([]uint64, (m+63)/64),
		keys: append([]uint64{}, keys...),
	}
}
