package set

import (
	"fmt"
)

type Copier interface {
	Copy() Copier
}

type Element interface {
	IsEqual(e Element) bool
}

type Newer interface {
	New() Newer
}

type Stringer interface {
	String() string
}

type Set interface {
	Copier
	Element
	Newer
	Stringer
	Add(elements ...Element)
	At(index int) (Element, bool)
	Probe(e Element) bool
	Remove(elements ...Element)
	Size() int
}

func Intersect(S, T Set) Set {
	ret := S.New().(Set)

	for i := 0; i < S.Size(); i += 1 {
		for j := 0; j < T.Size(); j += 1 {
			s, _ := S.At(i)
			t, _ := T.At(j)
			if s.IsEqual(t) {
				ret.Add(s)
				break
			}
		}
	}

	return ret
}

func Join(S, T Set) Set {
	ret := S.Copy().(Set)

	for i := 0; i < T.Size(); i += 1 {
		e, _ := T.At(i)
		ret.Add(e)
	}

	return ret
}

type simpleSet []Element

func NewSet() Set {
	return new(simpleSet)
}

func (S *simpleSet) Add(elements ...Element) {

	for _, e := range elements {
		if S.Probe(e) == false {
			*S = append(*S, e)
		}
	}
}

func (S simpleSet) At(index int) (Element, bool) {
	if index < len(S) {
		return S[index], true
	} //else {
	return S[0], false
	//}
}

func (S simpleSet) Copy() Copier {
	cp := make(simpleSet, S.Size())
	copy(cp, S)
	return &cp
}

func (S simpleSet) IsEqual(e Element) bool {
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

func (S simpleSet) New() Newer {
	return NewSet()
}

func (S simpleSet) Probe(e Element) bool {
	for _, v := range S {
		if e.IsEqual(v) {
			return true
		}
	}

	return false
}

func (S *simpleSet) Remove(elements ...Element) {

	for _, e := range elements {

		if S.Probe(e) == false {
			continue
		}

		cp := NewSet().(*simpleSet)

		for _, v := range *S {
			if e.IsEqual(v) == false {
				cp.Add(v)
			}
		}

		*S = *cp
	}
}

func (S simpleSet) Size() int {
	return len(S)
}

func (S simpleSet) String() string {
	return fmt.Sprintf("%v", []Element(S))
}
