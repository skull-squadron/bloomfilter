package bloomfilter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
)

const (
	keyFormat  = "%016x"
	bitsFormat = "%016x"
)

var nl = fmt.Sprintln()

func UnmarshalText(text []byte) (f *Filter, err error) {
	var k, n, m uint64
	r := bytes.NewBuffer(text)

	format := "k" + nl + "%d" + nl
	format += "n" + nl + "%d" + nl
	format += "m" + nl + "%d" + nl
	format += "keys" + nl

	_, err = fmt.Fscanf(r, format, k, n, m)
	if err != nil {
		return
	}

	if m < M_MIN {
		return nil, mError
	}

	if k < K_MIN {
		return nil, kError
	}

	f = &Filter{
		m:    m,
		n:    n,
		keys: make([]uint64, k),
		bits: make([]uint64, (m+63)/64),
	}

	for i := range f.keys {
		_, err = fmt.Fscanf(r, keyFormat, f.keys[i])
		if err != nil {
			return
		}
	}

	_, err = fmt.Fscanf(r, "bits")
	if err != nil {
		return
	}

	for i := range f.bits {
		_, err = fmt.Fscanf(r, bitsFormat, f.bits[i])
		if err != nil {
			return
		}
	}

	_, err = fmt.Fscanf(r, "sha384")
	if err != nil {
		return
	}

	actualHash := [sha512.Size384]byte{}

	for i := range actualHash {
		_, err = fmt.Fscanf(r, "%02x", actualHash[i])
		if err != nil {
			return
		}
	}

	_, expectedHash, err := f.marshal()
	if err != nil {
		return
	}

	if !hmac.Equal(expectedHash[:], actualHash[:]) {
		return nil, errHash
	}

	return
}

func (f *Filter) UnmarshalText(text []byte) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f2, err := UnmarshalText(text)
	if err != nil {
		return
	}

	f.m = f2.m
	f.n = f2.n
	copy(f.bits, f2.bits)
	copy(f.keys, f2.keys)

	return
}
