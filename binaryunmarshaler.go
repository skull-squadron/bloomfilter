// Package bloomfilter is face-meltingly fast, thread-safe,
// marshalable, unionable, probability- and
// optimal-size-calculating Bloom filter in go
//
// https://github.com/steakknife/bloomfilter
//
// Copyright © 2014, 2015, 2018 Barry Allard
//
// MIT license
//
package bloomfilter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
)

// UnmarshalBinary converts []bytes into a Filter
// conforms to encoding.BinaryUnmarshaler
func (f *Filter) UnmarshalBinary(data []byte) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	var k uint64

	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, &k)
	if err != nil {
		return err
	}

	if k < KMin {
		return errK
	}

	err = binary.Read(buf, binary.LittleEndian, &(f.n))
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian, &(f.m))
	if err != nil {
		return err
	}

	if f.m < MMin {
		return errM
	}

	debug("read bf k=%d n=%d m=%d\n", k, f.n, f.m)

	f.keys = make([]uint64, k)
	err = binary.Read(buf, binary.LittleEndian, f.keys)
	if err != nil {
		return err
	}

	f.bits = make([]uint64, f.n)
	err = binary.Read(buf, binary.LittleEndian, f.bits)
	if err != nil {
		return err
	}

	expectedHash := make([]byte, sha512.Size384)
	err = binary.Read(buf, binary.LittleEndian, expectedHash)
	if err != nil {
		return err
	}

	actualHash := sha512.Sum384(data[:len(data)-sha512.Size384])

	if !hmac.Equal(expectedHash, actualHash[:]) {
		debug("bloomfilter.UnmarshalBinary() sha384 hash failed:",
			" actual %v  expected %v", actualHash, expectedHash)
		return errHash
	}

	debug("bloomfilter.UnmarshalBinary() successfully read",
		" %d byte(s), sha384 %v", len(data), actualHash)
	return nil
}
