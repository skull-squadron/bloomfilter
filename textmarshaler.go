package bloomfilter

import (
	"fmt"
)

func (f *Filter) MarshalText() (text []byte, err error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	s := fmt.Sprintln("k")
	s += fmt.Sprintln(f.K())
	s += fmt.Sprintln("n")
	s += fmt.Sprintln(f.n)
	s += fmt.Sprintln("m")
	s += fmt.Sprintln(f.m)

	s += fmt.Sprintln("keys")
	for key := range f.keys {
		s += fmt.Sprintf(keyFormat, key) + nl
	}

	s += fmt.Sprintln("bits")
	for w := range f.bits {
		s += fmt.Sprintf(bitsFormat, w) + nl
	}

	_, hash, err := f.marshal()
	if err != nil {
		return
	}
	s += fmt.Sprintln("sha384")
	for b := range hash {
		s += fmt.Sprintf("%02x", b)
	}
	s += nl

	text = []byte(s)
	return
}
