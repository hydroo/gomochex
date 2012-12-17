package bitset

import (
	"fmt"
	"github.com/hydroo/gomochex/basic/set"
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

func NewBitSet() BitSet {
	return make(BitSet, 0)
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

	S.resizeIfPossible()

	// the bitset is completely empty
	*S = NewBitSet()
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

func (S *BitSet) resizeIfPossible() {
	for i := len(*S) - 1; i >= 0; i -= 1 {
		if (*S)[i] != 0x0 {
			if i < len(*S)-1 {
				S.resize(BitPosition(i * 64))
			}
			return
		}
	}
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

func Intersect(S, T BitSet) BitSet {
	U := NewBitSet()
	minLen := 0
	if len(S) > len(T) {
		minLen = len(T)
	} else {
		minLen = len(S)
	}

	if minLen == 0 {
		return U
	}

	U.resize(BitPosition((minLen * 64) - 1))

	for i := 0; i < minLen; i += 1 {
		U[i] = S[i] & T[i]
	}

	return U
}

func Join(S, T BitSet) BitSet {
	U := NewBitSet()

	maxLen := 0
	minLen := 0
	maxBitSet := S
	if len(S) > len(T) {
		maxLen = len(S)
		maxBitSet = S
		minLen = len(T)
	} else {
		maxLen = len(T)
		maxBitSet = T
		minLen = len(S)
	}

	if maxLen == 0 {
		return U
	}

	U.resize(BitPosition((maxLen * 64) - 1))

	for i := 0; i < minLen; i += 1 {
		U[i] = S[i] | T[i]
	}

	for i := minLen; i < maxLen; i += 1 {
		U[i] = maxBitSet[i]
	}

	return U
}
