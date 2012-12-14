package bitset_test

import (
	//"fmt"
	"github.com/hydroo/gomochex/basics/bitset"
	"testing"
)

func TestBitSetAdd(t *testing.T) {
	S := bitset.NewBitSet()
	S.Add(bitset.BitPosition(1), bitset.BitPosition(5), bitset.BitPosition(129))

	if S.Size() != 3 || len(*S) != 3 {
		t.Error()
	}

	s0, _ := S.At(0)
	s1, _ := S.At(1)
	s2, _ := S.At(2)
	if s0 != bitset.BitPosition(1) || s1 != bitset.BitPosition(5) || s2 != bitset.BitPosition(129) {
		t.Error()
	}
}

func TestBetSetAddDuplicates(t *testing.T) {
	S := bitset.NewBitSet()
	S.Add(bitset.BitPosition(1))
	S.Add(bitset.BitPosition(1))

	S.Add(bitset.BitPosition(1), bitset.BitPosition(1))

	if S.Size() != 1 || len(*S) != 1 {
		t.Error()
	}

	for i := 0; i < 64; i += 1 {
		if _, ok := S.At(i); ok != false && i != 0 {
			t.Error()
		}
	}
}

func TestBitSetIsEqual(t *testing.T) {
	S := bitset.NewBitSet()
	S.Add(bitset.BitPosition(1), bitset.BitPosition(2))

	T := bitset.NewBitSet()
	T.Add(bitset.BitPosition(1), bitset.BitPosition(2))

	if S.IsEqual(*T) != true {
		t.Error()
	}
}

func TestBitSetIsUnequal(t *testing.T) {
	S := bitset.NewBitSet()
	S.Add(bitset.BitPosition(1), bitset.BitPosition(2))

	T := bitset.NewBitSet()
	T.Add(bitset.BitPosition(1), bitset.BitPosition(3))


	U := bitset.NewBitSet()
	U.Add(bitset.BitPosition(1))

	V := bitset.BitPosition(1)

	if S.IsEqual(T) != false || U.IsEqual(V) != false || V.IsEqual(U) != false {
		t.Error()
	}
}

func TestBitSetRemove(t *testing.T) {
	S := bitset.NewBitSet()
	S.Add(bitset.BitPosition(0), bitset.BitPosition(1), bitset.BitPosition(2))

	S.Remove(bitset.BitPosition(1))

	if S.Size() != 2 || len(*S) != 1 {
		t.Error()
	}

	if S.Probe(bitset.BitPosition(0)) != true || S.Probe(bitset.BitPosition(1)) != false || S.Probe(bitset.BitPosition(2)) != true {
		t.Error()
	}
}

func TestBitSetRemoveFromEmpty(t *testing.T) {
	S := bitset.NewBitSet()

	S.Remove(bitset.BitPosition(1))

	if S.Size() != 0 || len(*S) != 0 {
		t.Error()
	}
}

func TestBitSetResize(t *testing.T) {
	S := bitset.NewBitSet()

	if S.Size() != 0 || len(*S) != 0 {
		t.Error()
	}

	S.Add(bitset.BitPosition(0))
	if S.Size() != 1 || len(*S) != 1 {
		t.Error()
	}

	S.Add(bitset.BitPosition(63))
	if S.Size() != 2 || len(*S) != 1 {
		t.Error()
	}

	S.Add(bitset.BitPosition(64))
	if S.Size() != 3 || len(*S) != 2 {
		t.Error()
	}

	S.Remove(bitset.BitPosition(63))
	if S.Size() != 2 || len(*S) != 2 {
		t.Error()
	}

	S.Remove(bitset.BitPosition(64))
	if S.Size() != 1 || len(*S) != 1 {
		t.Error()
	}

	S.Remove(bitset.BitPosition(0))
	if S.Size() != 0 || len(*S) != 0 {
		t.Error()
	}

}
