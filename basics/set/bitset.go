package set

import (
	"fmt"
)

type BitPosition int

func (b BitPosition) IsEqual(e Element) bool {
	return b == e.(BitPosition)
}

type BitSet []uint64

func NewBitSet() *BitSet {
	return new(BitSet)
}

func (S *BitSet) Add(elements ...Element) {

	for _, e := range elements {
		b := int(e.(BitPosition))
		if b >= len(*S)*64 {
			S.resize(BitPosition(b))
		}

		(*S)[b/64] = (*S)[b/64] | (1 << uint(b%64))
	}
}

func (S BitSet) At(index int) (Element, bool) {

	for i, v := range S {
		for j, b := 0, uint64(1<<0); j < 64; j, b = j+1, b<<1 {
			if v&b != 0 {
				if index == 0 {
					return BitPosition(i*64 + j), true
				}
				index -= 1
			}
		}
	}

	return BitPosition(-1), false
}

func (S BitSet) Copy() Copier {
	cp := NewBitSet()
	cp.resize(BitPosition((len(S) * 64) - 1))
	copy(*cp, S)
	return cp
}

func (S BitSet) IsEqual(e Element) bool {
	T := e.(Set)

	if S.Size() != T.Size() {
		return false
	}

	for i := 0; i < S.Size(); i += 1 {
		s, _ := S.At(i)

		if T.Probe(s) != true {
			return false
		}
	}

	return true
}

func (S BitSet) New() Newer {
	return new(BitSet)
}

func (S BitSet) Probe(e Element) bool {
	b := int(e.(BitPosition))
	if b >= len(S)*64 {
		return false
	}
	return (S[b/64] & (1 << uint(b%64))) != 0
}

func (S *BitSet) Remove(elements ...Element) {

	for _, e := range elements {
		b := int(e.(BitPosition))

		if b >= len(*S)*64 {
			continue
		}

		(*S)[b/64] = (*S)[b/64] & ^(1 << uint(b%64))
	}

	for i := len(*S) - 1; i >= 0; i -= 1 {
		if (*S)[i] != 0x0 {
			if i < len(*S)-1 {
				S.resize(BitPosition(i * 64))
			}
			return
		}
	}

	// the bitset is completely empty
	*S = *NewBitSet()
}

func (S *BitSet) resize(b BitPosition) {
	T := make(BitSet, int(b/64)+1)
	for k, v := range *S {
		if k >= len(T) {
			break
		}
		T[k] = v
	}
	*S = T
}

func (S BitSet) Size() int {

	size := 0

	for _, v := range S {
		for b := uint64(0x1); b != 0x0; b = b << 1 {
			if v&b != 0 {
				size += 1
			}
		}
	}

	return size
}

func (S BitSet) String() string {
	return fmt.Sprintf("%b", []uint64(S))
}
