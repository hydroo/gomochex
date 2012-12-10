package cartesianproduct

import (
	//"fmt"
	"github.com/hydroo/gomochex/basics/set"
)

type Tuple struct {
	first, second set.Element
}

func (t Tuple) IsEqual(u set.Element) bool {
	v := u.(Tuple)

	return t.First().IsEqual(v.First()) && t.Second().IsEqual(v.Second())
}

func (t Tuple) First() set.Element {
	return t.first
}

func (t Tuple) Second() set.Element {
	return t.second
}

func NewTuple(first, second set.Element) *Tuple {
	return &Tuple{first,second}
}


type CartesianProduct set.Set

func NewCartesianProduct(S, T set.Set) set.Set {

	SxT := set.NewSimpleSet()

	for i := 0; i < S.Size(); i += 1 {
		for j := 0; j < T.Size(); j += 1 {
			s, _ := S.At(i)
			t, _ := T.At(j)
			SxT.Add(*NewTuple(s, t))
		}
	}

	return SxT
}

