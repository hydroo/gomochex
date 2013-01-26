package bitset

import (
	//"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	S := NewBitSet()
	S.Add(1, 5, 129)

	if S.Size() != 3 || len(S) != 3 {
		t.Error()
	}

	s0, _ := S.At(0)
	s1, _ := S.At(1)
	s2, _ := S.At(2)
	if s0 != 1 || s1 != 5 || s2 != 129 {
		t.Error()
	}
}

func TestBetSetAddDuplicates(t *testing.T) {
	S := NewBitSet()
	S.Add(1)
	S.Add(1)

	S.Add(1, 1)

	if S.Size() != 1 || len(S) != 1 {
		t.Error()
	}

	for i := 0; i < 64; i += 1 {
		if _, ok := S.At(i); ok != false && i != 0 {
			t.Error()
		}
	}
}

func TestIsEqual(t *testing.T) {
	S := NewBitSet(1, 2)
	T := NewBitSet(1, 2)

	if S.IsEqual(T) != true {
		t.Error()
	}
}

func TestIsUnequal(t *testing.T) {
	S := NewBitSet(1, 2)
	T := NewBitSet(1, 3)
	U := NewBitSet(1)

	V := BitPosition(1)

	if S.IsEqual(T) != false || U.IsEqual(V) != false || V.IsEqual(U) != false {
		t.Error()
	}
}

func TestRemove(t *testing.T) {
	S := NewBitSet(0, 1, 2)

	S.Remove(1)

	if S.Size() != 2 || len(S) != 1 {
		t.Error()
	}

	if S.Probe(0) != true || S.Probe(1) != false || S.Probe(2) != true {
		t.Error()
	}
}

func TestRemoveFromEmpty(t *testing.T) {
	S := NewBitSet()

	S.Remove(1)

	if S.Size() != 0 || len(S) != 0 {
		t.Error()
	}
}

func TestResize(t *testing.T) {
	S := NewBitSet()

	if S.Size() != 0 || len(S) != 0 {
		t.Error()
	}

	S.Add(0)
	if S.Size() != 1 || len(S) != 1 {
		t.Error()
	}

	S.Add(63)
	if S.Size() != 2 || len(S) != 1 {
		t.Error()
	}

	S.Add(64)
	if S.Size() != 3 || len(S) != 2 {
		t.Error()
	}

	S.Remove(63)
	if S.Size() != 2 || len(S) != 2 {
		t.Error()
	}

	S.Remove(64)
	if S.Size() != 1 || len(S) != 1 {
		t.Error()
	}

	S.Remove(0)
	if S.Size() != 0 || len(S) != 0 {
		t.Error()
	}

}

func TestIntersect(t *testing.T) {

	S := NewBitSet(1, 2, 3)
	T := NewBitSet(2, 3, 64)
	U := Intersect(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	if U.Size() != 2 || u0 != 2 || u1 != 3 {
		t.Error()
	}
}

func TestJoin(t *testing.T) {

	S := NewBitSet(1, 2, 3)
	T := NewBitSet(2, 3, 64)
	U := Join(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	u2, _ := U.At(2)
	u3, _ := U.At(3)

	if U.Size() != 4 || u0 != 1 || u1 != 2 || u2 != 3 || u3 != 64 {
		t.Error()
	}
}
