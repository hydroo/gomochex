package cartesianproduct_test

import (
	//"fmt"
	"github.com/hydroo/gomochex/basics/cartesianproduct"
	"github.com/hydroo/gomochex/basics/set"
	"testing"
)

type myInt int

func (i myInt) IsEqual(j set.Element) bool {
	if k, ok := j.(myInt); ok == true && i == k {
		return true
	} //else {
	return false
	//}
}

func TestNewCartesianProduct(t *testing.T) {
	S := set.NewSet()
	S.Add(myInt(0), myInt(1), myInt(2))

	U := set.NewSet()
	U.Add(myInt(3), myInt(4), myInt(5))

	SxU := cartesianproduct.NewCartesianProduct(S, U)

	if SxU.Size() != 9 {
		t.Error()
	}

	for i := 0; i < SxU.Size(); i += 1 {
		w, ok := SxU.At(i)
		v := w.(cartesianproduct.Tuple)

		s, _ := S.At(i / 3)
		u, _ := U.At(i % 3)

		if ok != true || v.First() != s || v.Second() != u {
			t.Error()
		}
	}
}
