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
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

// ReadFrom r and replace f with the new bf
func (f *Filter) ReadFrom(r io.Reader) (n int64, err error) {
	f2, n, err := ReadFrom(r)
	if err != nil {
		return -1, err
	}
	*f = *f2
	return n, nil
}

// Read a lossless-compressed Bloom filter from Reader
func ReadFrom(r io.Reader) (f *Filter, n int64, err error) {
	rawR, err := gzip.NewReader(r)
	if err != nil {
		return nil, -1, err
	}
	defer func() {
		err = rawR.Close()
	}()

	content, err := ioutil.ReadAll(rawR)
	if err != nil {
		return nil, -1, err
	}

	f = new(Filter)
	n = int64(len(content))
	err = f.UnmarshalBinary(content)
	if err != nil {
		return nil, -1, err
	}
	return f, n, nil
}

// Read a lossless-compressed Bloom filter to a file
// Suggested file extension: .bf.gz
func ReadFile(filename string) (_ *Filter, n int64, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, -1, err
	}
	defer func() {
		err = r.Close()
	}()

	return ReadFrom(r)
}

// Write a lossless-compressed Bloom filter to Writer
func (f *Filter) WriteTo(w io.Writer) (n int64, err error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	rawW := gzip.NewWriter(w)
	defer func() {
		err = rawW.Close()
	}()

	content, err := f.MarshalBinary()
	if err != nil {
		return -1, err
	}

	intN, err := rawW.Write(content)
	n = int64(intN)
	return n, err
}

// Write a lossless-compressed Bloom filter to a file
// Suggested file extension: .bf.gz
func (f *Filter) WriteFile(filename string) (n int64, err error) {
	w, err := os.Create(filename)
	if err != nil {
		return -1, err
	}
	defer func() {
		err = w.Close()
	}()

	return f.WriteTo(w)
}
