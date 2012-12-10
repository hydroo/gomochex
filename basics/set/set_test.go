package set_test

import (
	//"fmt"
	"github.com/hydroo/gomochex/basics/set"
	"testing"
)

type myInt int

func (i myInt) IsEqual(j set.Element) bool {
	return i == j.(myInt)
}

func TestSetInit(t *testing.T) {
	S := set.NewSimpleSet()

	if S.Size() != 0 {
		t.Error()
	}
}

func TestSetAdd(t *testing.T) {
	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(2))
	S.Add(myInt(3))

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

func TestSetAddDuplicates(t *testing.T) {
	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(1))

	if S.Size() != 1 {
		t.Error()
	}
}

func TestSetIsEqual(t *testing.T) {
	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(2))

	T := set.NewSimpleSet()
	T.Add(myInt(1))
	T.Add(myInt(2))

	if S.IsEqual(T) != true {
		t.Error()
	}
}

func TestSetIsUnequal(t *testing.T) {
	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(2))

	T := set.NewSimpleSet()
	T.Add(myInt(1))
	T.Add(myInt(3))

	if S.IsEqual(T) != false {
		t.Error()
	}
}

func TestSetRemove(t *testing.T) {
	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(2))
	S.Add(myInt(3))

	S.Remove(myInt(2))

	if S.Size() != 2 {
		t.Error()
	}

	if S.Probe(myInt(1)) != true || S.Probe(myInt(2)) != false || S.Probe(myInt(3)) != true {
		t.Error()
	}
}

func TestSetRemoveFromEmpty(t *testing.T) {
	S := set.NewSimpleSet()

	S.Remove(myInt(1))

	if S.Size() != 0 {
		t.Error()
	}
}

func TestSetIntersect(t *testing.T) {

	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(2))
	S.Add(myInt(3))

	T := set.NewSimpleSet()
	T.Add(myInt(2))
	T.Add(myInt(3))
	T.Add(myInt(4))

	U := set.Intersect(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	if U.Size() != 2 || u0 != myInt(2) || u1 != myInt(3) {
		t.Error()
	}
}

func TestSetJoin(t *testing.T) {

	S := set.NewSimpleSet()
	S.Add(myInt(1))
	S.Add(myInt(2))
	S.Add(myInt(3))

	T := set.NewSimpleSet()
	T.Add(myInt(2))
	T.Add(myInt(3))
	T.Add(myInt(4))

	U := set.Join(S, T)

	u0, _ := U.At(0)
	u1, _ := U.At(1)
	u2, _ := U.At(2)
	u3, _ := U.At(3)

	if U.Size() != 4 || u0 != myInt(1) || u1 != myInt(2) || u2 != myInt(3) || u3 != myInt(4) {
		t.Error()
	}
}
