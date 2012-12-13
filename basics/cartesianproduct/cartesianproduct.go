package cartesianproduct

import (
	//"fmt"
	"github.com/hydroo/gomochex/basics/set"
)

type Tuple interface {
	set.Element
	First() set.Element
	Second() set.Element
}

type simpleTuple struct {
	first, second set.Element
}

func (t simpleTuple) IsEqual(u set.Element) bool {
	if v, ok := u.(Tuple); ok == true && t.First().IsEqual(v.First()) && t.Second().IsEqual(v.Second()) {
		return true
	} //else {
	return false
	//}
}

func (t simpleTuple) First() set.Element {
	return t.first
}

func (t simpleTuple) Second() set.Element {
	return t.second
}

func NewTuple(first, second set.Element) Tuple {
	return &simpleTuple{first, second}
}


func NewCartesianProduct(S, T set.Set) set.Set {

	SxT := set.NewSet()

	for i := 0; i < S.Size(); i += 1 {
		for j := 0; j < T.Size(); j += 1 {
			s, _ := S.At(i)
			t, _ := T.At(j)
			SxT.Add(NewTuple(s, t))
		}
	}

	return SxT
}

