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

// conforms to encoding.BinaryUnmarshaler

func (f *Filter) UnmarshalBinary(data []byte) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	var k uint32

	buf := bytes.NewBuffer(data)
	if err = binary.Read(buf, binary.LittleEndian, &k); err != nil {
		return err
	}

	if err = binary.Read(buf, binary.LittleEndian, &(f.n)); err != nil {
		return err
	}

	if err = binary.Read(buf, binary.LittleEndian, &(f.m)); err != nil {
		return err
	}

	debug("read bf k=%d n=%d m=%d\n", k, f.n, f.m)

	f.keys = make([]uint64, k)
	if err = binary.Read(buf, binary.LittleEndian, f.keys); err != nil {
		return err
	}

	f.bits = make([]uint64, f.n)
	if err = binary.Read(buf, binary.LittleEndian, f.bits); err != nil {
		return err
	}

	return nil
}
