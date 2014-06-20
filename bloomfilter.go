package bloomfilter

import (
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

func countOnes(x uint64) int {
  c := x - ((x >> 1) & 033333333333) - ((x >> 2) & 011111111111)
  return int((c+(c>>3))&030707070707) % 63
}

// count # of 1's
func (f Filter) FilledRatio() float32 {
  ones := 0
  for _, b := range f.bits {
    ones += countOnes(b)
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
