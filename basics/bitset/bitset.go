package bitset

import (
	"fmt"
	"github.com/hydroo/gomochex/basics/set"
)

type BitPosition int

func (b BitPosition) IsEqual(e set.Element) bool {
	if f, ok := e.(BitPosition); ok == true && b == f {
		return true
	} //else {
	return false
	//}
}

type BitSet []uint64

func NewBitSet() *BitSet {
	return &BitSet{}
}

func (S *BitSet) Add(elements ...BitPosition) {

	for _, e := range elements {
		b := int(e)
		if b >= len(*S)*64 {
			S.resize(BitPosition(b))
		}

		(*S)[b/64] = (*S)[b/64] | (1 << uint(b%64))
	}
}

func (S BitSet) At(index int) (BitPosition, bool) {

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

func (S BitSet) Copy() set.Copier {
	cp := NewBitSet()
	cp.resize(BitPosition((len(S) * 64) - 1))
	copy(*cp, S)
	return cp
}

func (S BitSet) IsEqual(e set.Element) bool {
	T, ok := e.(BitSet)

	if ok == false {
		return false
	}

	if len(S) != len(T) {
		return false
	}

	for i := 0; i < len(S); i += 1 {
		if S[i] != T[i] {
			return false
		}
	}

	return true
}

func (S BitSet) New() set.Newer {
	return NewBitSet()
}

func (S BitSet) Probe(b BitPosition) bool {
	if int(b) >= len(S)*64 {
		return false
	}
	return (S[b/64] & (1 << uint(b%64))) != 0
}

func (S *BitSet) Remove(elements ...BitPosition) {

	for _, e := range elements {
		b := int(e)

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
