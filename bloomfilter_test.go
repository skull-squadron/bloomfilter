package bloomfilter

import (
  "math/rand"
  "testing"
)

type HashableUint32 uint32

func (h HashableUint32) BloomFilterHash() uint32 {
  return uint32(h)
}

var hashableUint32Values = []HashableUint32{
  0,
  7,
  0x0c0ffee0,
  0xdeadbeef,
  0xffffffff,
}

var hashableUint32NotValues = []HashableUint32{
  1,
  5,
  42,
  0xa5a5a5a5,
  0xfffffffe,
}

func Test0(t *testing.T) {
  bf := New(10000, 5)

  t.Log("Filled ratio before adds :", bf.FilledRatio())
  for _, x := range hashableUint32Values {
    bf.Add(x)
  }
  t.Log("Filled ratio after adds :", bf.FilledRatio())

  // these may or may not be true
  for _, y := range hashableUint32Values {
    if bf.Contains(y) {
      t.Log("value in set querties: may contain ", y)
    } else {
      t.Fatal("value in set queries: definitely does not contain ", y, ", but it should")
    }
  }

  // these must all be false
  for _, z := range hashableUint32NotValues {
    if bf.Contains(z) {
      t.Log("value not in set queries: may or may not contain ", z)
    } else {
      t.Log("value not in set queries: definitely does not contain ", z, " which is correct")
    }
  }
}

func BenchmarkAddX10kX5(b *testing.B) {
  b.StopTimer()
  bf := New(10000, 5)
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    bf.Add(HashableUint32(rand.Uint32()))
  }
}

func BenchmarkContains1kX10kX5(b *testing.B) {
  b.StopTimer()
  bf := New(10000, 5)
  for i := 0; i < 1000; i++ {
    bf.Add(HashableUint32(rand.Uint32()))
  }
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    bf.Contains(HashableUint32(rand.Uint32()))
  }
}

func BenchmarkContains100kX10BX20(b *testing.B) {
  b.StopTimer()
  bf := New(10*1000*1000*1000, 20)
  for i := 0; i < 100*1000; i++ {
    bf.Add(HashableUint32(rand.Uint32()))
  }
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    bf.Contains(HashableUint32(rand.Uint32()))
  }
}
