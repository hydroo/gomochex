package regex_test

import (
	"fmt"
	"github.com/hydroo/gomochex/automaton/nfa"
	"github.com/hydroo/gomochex/regex"
	"testing"
)

func TestExpressionFromString(t *testing.T) {

	//correct concatenation
	if e, ok := regex.ExpressionFromString("(asdf.π)"); ok != true || fmt.Sprint(e) != "(asdf.π)" {
		t.Error()
	}

	//correct letter
	if e, ok := regex.ExpressionFromString("π"); ok != true || fmt.Sprint(e) != "π" {
		t.Error()
	}

	//wrong letter
	if _, ok := regex.ExpressionFromString("π."); ok != false {
		t.Error()
	}

	//correct letter
	if e, ok := regex.ExpressionFromString("πasdf"); ok != true || fmt.Sprint(e) != "πasdf" {
		t.Error()
	}

	//correct or
	if e, ok := regex.ExpressionFromString("(asdf+π)"); ok != true || fmt.Sprint(e) != "(asdf+π)" {
		t.Error()
	}

	//wrong or
	if _, ok := regex.ExpressionFromString("((asdf+π)"); ok != false {
		t.Error()
	}

	//correct star
	if e, ok := regex.ExpressionFromString("((asdf+π))*"); ok != true || fmt.Sprint(e) != "((asdf+π))*" {
		t.Error()
	}

	//wrong star
	if _, ok := regex.ExpressionFromString("π*"); ok != false {
		t.Error()
	}

	//correct complex expr
	if e, ok := regex.ExpressionFromString("(a.((π+b).(c)*))"); ok != true || fmt.Sprint(e) != "(a.((π+b).(c)*))" {
		t.Error()
	}

}

func TestConNfa(t *testing.T) {
	//TODO
}

func TestLetterNfa(t *testing.T) {

	a := nfa.StringLetter("a")

	A := regex.Letter("a").Nfa()

	if A.Alphabet().Size() != 1 || A.States().Size() != 2 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}

	if b, ok := A.Alphabet().At(0); ok != true || b.IsEqual(nfa.StringLetter(a)) == false {
		t.Error()
	}

	//has exactly one transition which is not a loop,
	//and goes from an initial to a final state
	for i := 0; i < A.States().Size(); i += 1 {
		for j := 0; j < A.Alphabet().Size(); j += 1 {

			q0, _ := A.States().At(i)
			b, _ := A.Alphabet().At(j)

			S := A.Transition(q0, b)

			if S.Size() != 0 {

				if S.Size() != 1 {
					t.Error()
				}

				qf, _ := S.At(0)

				if q0.IsEqual(qf) || A.InitialStates().Probe(q0) == false || A.FinalStates().Probe(qf) == false {
					t.Error()
				}
			}

		}
	}
}

func TestOrNfa(t *testing.T) {
	//TODO
}

func TestStarNfa(t *testing.T) {
	//TODO
}
