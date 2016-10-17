package bloomfilter

import (
	"bytes"
	"testing"
)

func TestWriteRead(t *testing.T) {
	// minimal filter
	f := New(2, 1)

	v := hashableUint64(0)
	f.Add(v)

	var b bytes.Buffer

	_, err := f.WriteTo(&b)
	if err != nil {
		t.Error(err)
	}

	f2, _, err := ReadFrom(&b)
	if err != nil {
		t.Error(err)
	}

	if !f2.Contains(v) {
		t.Error("Filters not equal")
	}
}
