package set_test

import (
	//"fmt"
	"github.com/hydroo/gomochex/basics/set"
	"testing"
)

func TestBitSetAdd(t *testing.T) {
	S := set.NewBitSet()
	S.Add(set.BitPosition(1))
	S.Add(set.BitPosition(5))
	S.Add(set.BitPosition(129))

	if S.Size() != 3 || len(*S) != 3 {
		t.Error()
	}

	s0, _ := S.At(0)
	s1, _ := S.At(1)
	s2, _ := S.At(2)
	if s0 != set.BitPosition(1) || s1 != set.BitPosition(5) || s2 != set.BitPosition(129) {
		t.Error()
	}
}

func TestBetSetAddDuplicates(t *testing.T) {
	S := set.NewBitSet()
	S.Add(set.BitPosition(1))
	S.Add(set.BitPosition(1))

	if S.Size() != 1 || len(*S) != 1 {
		t.Error()
	}

	for i := 0; i < 64; i += 1 {
		if _, ok := S.At(i); ok != false && i != 0 {
			t.Error()
		}
	}
}

func TestBitSetRemove(t *testing.T) {
	S := set.NewBitSet()
	S.Add(set.BitPosition(0))
	S.Add(set.BitPosition(1))
	S.Add(set.BitPosition(2))

	S.Remove(set.BitPosition(1))

	if S.Size() != 2 || len(*S) != 1 {
		t.Error()
	}

	if S.Probe(set.BitPosition(0)) != true || S.Probe(set.BitPosition(1)) != false || S.Probe(set.BitPosition(2)) != true {
		t.Error()
	}
}

func TestBitSetRemoveFromEmpty(t *testing.T) {
	S := set.NewBitSet()

	S.Remove(set.BitPosition(1))

	if S.Size() != 0 || len(*S) != 0 {
		t.Error()
	}
}

func TestBitSetResize(t *testing.T) {
	S := set.NewBitSet()

	if S.Size() != 0 || len(*S) != 0 {
		t.Error()
	}

	S.Add(set.BitPosition(0))
	if S.Size() != 1 || len(*S) != 1 {
		t.Error()
	}

	S.Add(set.BitPosition(63))
	if S.Size() != 2 || len(*S) != 1 {
		t.Error()
	}

	S.Add(set.BitPosition(64))
	if S.Size() != 3 || len(*S) != 2 {
		t.Error()
	}

	S.Remove(set.BitPosition(63))
	if S.Size() != 2 || len(*S) != 2 {
		t.Error()
	}

	S.Remove(set.BitPosition(64))
	if S.Size() != 1 || len(*S) != 1 {
		t.Error()
	}

	S.Remove(set.BitPosition(0))
	if S.Size() != 0 || len(*S) != 0 {
		t.Error()
	}

}

func TestBitSetIntersect(t *testing.T) {

	S := new(set.BitSet)
	S.Add(set.BitPosition(0))
	S.Add(set.BitPosition(1))
	S.Add(set.BitPosition(2))

	T := new(set.BitSet)
	T.Add(set.BitPosition(1))
	T.Add(set.BitPosition(2))
	T.Add(set.BitPosition(3))

	U := set.Intersect(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	if U.Size() != 2 || u0 != set.BitPosition(1) || u1 != set.BitPosition(2) {
		t.Error()
	}
}

func TestBitSetJoin(t *testing.T) {

	S := new(set.BitSet)
	S.Add(set.BitPosition(0))
	S.Add(set.BitPosition(1))
	S.Add(set.BitPosition(2))

	T := new(set.BitSet)
	T.Add(set.BitPosition(1))
	T.Add(set.BitPosition(2))
	T.Add(set.BitPosition(3))

	U := set.Join(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	u2, _ := U.At(2)
	u3, _ := U.At(3)
	if U.Size() != 4 || u0 != set.BitPosition(0) || u1 != set.BitPosition(1) || u2 != set.BitPosition(2) || u3 != set.BitPosition(3) {
		t.Error()
	}
}
