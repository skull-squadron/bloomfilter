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
	"bytes"
	"encoding/binary"
)

// conforms to encoding.BinaryMarshaler

// marshalled binary layout:
//
//   k
//   n
//   m
//   keys
//   bits
//
func (f *Filter) MarshalBinary() (data []byte, err error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, f.K()); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.LittleEndian, f.n); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.LittleEndian, f.m); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.LittleEndian, f.keys); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.LittleEndian, f.bits); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
