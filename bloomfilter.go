package bloomfilter

import (
  "github.com/steakknife/hamming"
  "math/rand"
  "time"
)

const (
  randSeedMagic = 0x4ae5
)

type Hashable interface {
  BloomFilterHash() uint32
}

type Filter struct {
  bits []uint64
  keys []uint32
}

func newKeys(k int) (keys []uint32) {
  keys = make([]uint32, k)
  r := rand.New(rand.NewSource(time.Now().UnixNano() ^ randSeedMagic))
  for i := 0; i < k; i++ {
    keys[i] = r.Uint32()
  }
  return
}

func New(m int, k int) *Filter {
  return &Filter{
    bits: make([]uint64, (m+63)/64),
    keys: newKeys(k),
  }
}

func (f Filter) hashKeys(v Hashable) (hashes []uint32) {
  rawHash := v.BloomFilterHash()
  n := len(f.keys)
  hashes = make([]uint32, n)
  for i := 0; i < n; i++ {
    hashes[i] = rawHash ^ f.keys[i]
  }
  return
}

func (f Filter) getBit(i uint32) bool {
  n := uint32(len(f.bits))
  if i >= n {
    i %= n
  }
  return (f.bits[i>>6] >> uint(i&0x3f)) != 0
}

func (f Filter) setBit(i uint32) {
  n := uint32(len(f.bits))
  if i >= n {
    i %= n
  }
  f.bits[i>>6] |= 1 << uint(i&0x3f)
}

func (f *Filter) Add(v Hashable) {
  for _, k := range f.hashKeys(v) {
    f.setBit(k)
  }
}

/*
1  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3  3
1011011011011011011011011011011011011011011011011011011011011011

1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1  1
1001001001001001001001001001001001001001001001001001001001001001

0  7  0  7  0  7  0  7  0  7  0  7  0  7  0  7  0  7  0  7  0  7
0111000111000111000111000111000111000111000111000111000111000111

0123456789012345678901234567890123456789012345678901234567890123
          1         2         3         4         5         6
*/

const (
  m1  uint64 = 0x5555555555555555 //binary: 0101...
  m2  uint64 = 0x3333333333333333 //binary: 00110011..
  m4  uint64 = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...
  m8  uint64 = 0x00ff00ff00ff00ff //binary:  8 zeros,  8 ones ...
  m16 uint64 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
  m32 uint64 = 0x00000000ffffffff //binary: 32 zeros, 32 ones
  hff uint64 = 0xffffffffffffffff //binary: all ones
  h01 uint64 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...
)

// count # of 1's
func (f Filter) FilledRatio() float32 {
  ones := 0
  for _, b := range f.bits {
    ones += hamming.CountBitsUint64(b)
  }
  return float32(ones) / float32(len(f.bits))
}

// false: definitely false
// true:  maybe true or false
func (f Filter) Contains(v Hashable) bool {
  for _, k := range f.hashKeys(v) {
    if !f.getBit(k) {
      return false
    }
  }
  return true // maybe
}
