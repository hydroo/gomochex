package set

import (
	//"fmt"
	"testing"
)

type myInt int

func (i myInt) IsEqual(j Element) bool {

	if k, ok := j.(myInt); ok == true && i == k {
		return true
	} //else {
	return false
	//}
}

func TestInit(t *testing.T) {
	S := NewSet()

	if S.Size() != 0 {
		t.Error()
	}
}

func TestAdd(t *testing.T) {
	S := NewSet()
	S.Add(myInt(1), myInt(2), myInt(3))

	if S.Size() != 3 {
		t.Error()
	}

	s0, _ := S.At(0)
	s1, _ := S.At(1)
	s2, _ := S.At(2)
	if s0 != myInt(1) || s1 != myInt(2) || s2 != myInt(3) {
		t.Error()
	}
}

func TestAddDuplicates(t *testing.T) {
	S := NewSet()
	S.Add(myInt(1))
	S.Add(myInt(1))

	S.Add(myInt(1), myInt(1))

	if S.Size() != 1 {
		t.Error()
	}
}

func TestIsEqual(t *testing.T) {
	S := NewSet(myInt(1), myInt(2))

	T := NewSet(myInt(1), myInt(2))

	if S.IsEqual(T) != true {
		t.Error()
	}
}

func TestIsUnequal(t *testing.T) {
	S := NewSet(myInt(1), myInt(2))

	T := NewSet(myInt(1), myInt(3))

	U := NewSet(myInt(1))

	V := myInt(1)

	if S.IsEqual(T) != false || U.IsEqual(V) != false || V.IsEqual(U) != false {
		t.Error()
	}
}

func TestRemove(t *testing.T) {
	S := NewSet(myInt(1), myInt(2), myInt(3))

	S.Remove(myInt(2))

	if S.Size() != 2 {
		t.Error()
	}

	if S.Probe(myInt(1)) != true || S.Probe(myInt(2)) != false || S.Probe(myInt(3)) != true {
		t.Error()
	}
}

func TestRemoveFromEmpty(t *testing.T) {
	S := NewSet()

	S.Remove(myInt(1))

	if S.Size() != 0 {
		t.Error()
	}
}

func TestIntersect(t *testing.T) {

	S := NewSet(myInt(1), myInt(2), myInt(3))
	T := NewSet(myInt(2), myInt(3), myInt(4))
	U := Intersect(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	if U.Size() != 2 || u0 != myInt(2) || u1 != myInt(3) {
		t.Error()
	}
}

func TestJoin(t *testing.T) {

	S := NewSet(myInt(1), myInt(2), myInt(3))
	T := NewSet(myInt(2), myInt(3), myInt(4))
	U := Join(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	u2, _ := U.At(2)
	u3, _ := U.At(3)

	if U.Size() != 4 || u0 != myInt(1) || u1 != myInt(2) || u2 != myInt(3) || u3 != myInt(4) {
		t.Error()
	}
}
