package bloomfilter

// TODO probabilities, filled ratio

import (
  "bytes"
  "encoding/binary"
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

// m is the size of the bloom filter in bits
// k is the number of randomly generated keys used (One-Time-Pad-inspired)
// m must be > 0, or a runtime panic occurs
// k must be > 0, or a runtime panic occurs
func New(m int64, k int) *Filter {
  if m <= 0 {
    panic("m (number of bits in the bloom filter) must be > 0")
  }
  if k <= 0 {
    panic("k (number of keys uses in the bloom filter) must be > 0")
  }
  return &Filter{
    bits: make([]uint64, (m+63)>>6),
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

// 32LE(len(f.keys))
// 64LE(len(f.bits))
func (f Filter) MarshalBinary() (data []byte, err error) {
  nKeys := uint32(len(f.keys))
  nBits := uint64(len(f.bits))

  size := binary.Size(nKeys) + binary.Size(nBits) + binary.Size(f.keys) + binary.Size(f.bits)
  data = make([]byte, 0, size)
  buf := bytes.NewBuffer(data)

  err = binary.Write(buf, binary.LittleEndian, nKeys)
  if err != nil {
    return
  }

  err = binary.Write(buf, binary.LittleEndian, nBits)
  if err != nil {
    return
  }

  err = binary.Write(buf, binary.LittleEndian, f.keys)
  if err != nil {
    return
  }

  err = binary.Write(buf, binary.LittleEndian, f.bits)
  if err != nil {
    return
  }

  data = buf.Bytes()
  return
}

func (f *Filter) UnmarshalBinary(data []byte) (err error) {
  var (
    nKeys uint32
    nBits uint64
  )

  buf := bytes.NewBuffer(data)
  err = binary.Read(buf, binary.LittleEndian, nKeys)
  if err != nil {
    return
  }

  err = binary.Read(buf, binary.LittleEndian, nBits)
  if err != nil {
    return
  }

  f.keys = make([]uint32, nKeys, nKeys)
  err = binary.Read(buf, binary.LittleEndian, f.keys)
  if err != nil {
    return
  }

  f.bits = make([]uint64, nBits, nBits)
  err = binary.Read(buf, binary.LittleEndian, f.bits)
  if err != nil {
    return
  }

  return nil
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
