package cartesianproduct_test

import (
	//"fmt"
	"testing"
	"github.com/hydroo/gomochex/basics/set"
	"github.com/hydroo/gomochex/basics/cartesianproduct"
)


type myInt int

func (i myInt) IsEqual(j set.Element) bool {
	return i == j.(myInt)
}


func TestNewCartesianProduct(t *testing.T) {
	S := set.NewSimpleSet()
	S.Add(myInt(0))
	S.Add(myInt(1))
	S.Add(myInt(2))
	
	U := set.NewSimpleSet()
	U.Add(myInt(3))
	U.Add(myInt(4))
	U.Add(myInt(5))

	SxU := cartesianproduct.NewCartesianProduct(S, U)

	if SxU.Size() != 9 {
		t.Error()
	}

	for i := 0; i < SxU.Size(); i += 1 {
		w, ok := SxU.At(i)
		v := w.(cartesianproduct.Tuple)

		s ,_ := S.At(i/3)
		u ,_ := U.At(i%3)

		if ok != true || v.First() != s || v.Second() != u {
			t.Error()
		}
	}
}

