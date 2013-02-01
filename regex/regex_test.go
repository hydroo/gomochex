package regex

import (
	//"fmt"
	"github.com/hydroo/gomochex/automaton/nfa"
	"github.com/hydroo/gomochex/basic/set"
	"testing"
)

func TestExpressionFromString(t *testing.T) {

	//correct concatenation
	if e, ok := ExpressionFromString("(asdf.π)"); ok != true || e.IsEqual(Concat(Letter("asdf"), Letter("π"))) != true {
		t.Error()
	}

	//correct letter
	if e, ok := ExpressionFromString("π"); ok != true || e.IsEqual(Letter("π")) != true {
		t.Error()
	}

	//wrong letter
	if _, ok := ExpressionFromString("π."); ok != false {
		t.Error()
	}

	//correct letter
	if e, ok := ExpressionFromString("πasdf"); ok != true || e.IsEqual(Letter("πasdf")) != true {
		t.Error()
	}

	//correct or
	if e, ok := ExpressionFromString("(asdf+π)"); ok != true || e.IsEqual(Or(Letter("asdf"), Letter("π"))) != true {
		t.Error()
	}

	//wrong or
	if _, ok := ExpressionFromString("((asdf+π)"); ok != false {
		t.Error()
	}

	//correct star
	if e, ok := ExpressionFromString("((asdf+π))*"); ok != true || e.IsEqual(Star(Or(Letter("asdf"), Letter("π")))) != true {
		t.Error()
	}

	//wrong star
	if _, ok := ExpressionFromString("π*"); ok != false {
		t.Error()
	}

	//correct complex expr
	if e, ok := ExpressionFromString("(a.((π+b).(c)*))"); ok != true || e.IsEqual(Concat(Letter("a"), Concat(Or(Letter("π"), Letter("b")), Star(Letter("c"))))) != true {
		t.Error()
	}
}
