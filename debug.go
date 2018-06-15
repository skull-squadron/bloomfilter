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

import "log"

// Debug controls whether debug() logs anything
var Debug = false

// debug printing when Debug is true
func debug(format string, a ...interface{}) {
	if Debug {
		log.Printf(format, a...)
	}
}
