package regex_test

import (
	//"fmt"
	"github.com/hydroo/gomochex/automaton/nfa"
	"github.com/hydroo/gomochex/basic/set"
	"github.com/hydroo/gomochex/regex"
	"testing"
)

func TestExpressionFromString(t *testing.T) {

	//correct concatenation
	if e, ok := regex.ExpressionFromString("(asdf.π)"); ok != true || e.IsEqual(regex.Concat(regex.Letter("asdf"), regex.Letter("π"))) != true {
		t.Error()
	}

	//correct letter
	if e, ok := regex.ExpressionFromString("π"); ok != true || e.IsEqual(regex.Letter("π")) != true {
		t.Error()
	}

	//wrong letter
	if _, ok := regex.ExpressionFromString("π."); ok != false {
		t.Error()
	}

	//correct letter
	if e, ok := regex.ExpressionFromString("πasdf"); ok != true || e.IsEqual(regex.Letter("πasdf")) != true {
		t.Error()
	}

	//correct or
	if e, ok := regex.ExpressionFromString("(asdf+π)"); ok != true || e.IsEqual(regex.Or(regex.Letter("asdf"), regex.Letter("π"))) != true {
		t.Error()
	}

	//wrong or
	if _, ok := regex.ExpressionFromString("((asdf+π)"); ok != false {
		t.Error()
	}

	//correct star
	if e, ok := regex.ExpressionFromString("((asdf+π))*"); ok != true || e.IsEqual(regex.Star(regex.Or(regex.Letter("asdf"), regex.Letter("π")))) != true {
		t.Error()
	}

	//wrong star
	if _, ok := regex.ExpressionFromString("π*"); ok != false {
		t.Error()
	}

	//correct complex expr
	if e, ok := regex.ExpressionFromString("(a.((π+b).(c)*))"); ok != true || e.IsEqual(regex.Concat(regex.Letter("a"), regex.Concat(regex.Or(regex.Letter("π"), regex.Letter("b")), regex.Star(regex.Letter("c"))))) != true {
		t.Error()
	}

}

func TestConcatNfa(t *testing.T) {
	a := nfa.Letter("a")
	b := nfa.Letter("π")
	c := nfa.Letter("c")

	A := regex.Concat(regex.Letter("a"), regex.Or(regex.Letter("π"), regex.Letter("c"))).Nfa()

	if A.Alphabet().Size() != 3 || A.States().Size() < 5 || A.States().Size() > 6 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 2 {
		t.Error()
	}
	if A.Alphabet().Probe(a) != true || A.Alphabet().Probe(b) != true || A.Alphabet().Probe(c) != true {
		t.Error()
	}

	//       -- a --> o -- π --> □
	//      /
	// --> o
	//      \
	//       -- a --> o -- c --> □
	transitionCount := 0
	var lastLetter set.Element = nil //makes sure no two π or c transitions occur
	visitedRStates := set.NewSet()   //makes sure the two π and c transitions use entirely different states
	for i := 0; i < A.States().Size(); i += 1 {
		s, _ := A.States().At(i)
		for j := 0; j < A.Alphabet().Size(); j += 1 {
			x, _ := A.Alphabet().At(j)
			S := A.Transition(s.(nfa.State), x.(nfa.Letter))

			if S.Size() == 0 {
				continue
			}
			if S.Size() == 2 { // init -- a --> 2x
				u, _ := S.At(0)
				v, _ := S.At(1)
				if A.InitialStates().Probe(s) != true || A.InitialStates().Probe(u) != false || A.InitialStates().Probe(v) != false || A.FinalStates().Probe(s) != false || A.FinalStates().Probe(u) != false || A.FinalStates().Probe(v) != false || x.IsEqual(a) != true {
					t.Error()
				}
			} else if S.Size() == 1 { // o -- c --> □ or o -- π --> □
				u, _ := S.At(0)

				if x.IsEqual(b) == false && x.IsEqual(c) == false {
					t.Error()
				}

				if A.InitialStates().Probe(s) != false || A.FinalStates().Probe(s) != false || A.InitialStates().Probe(u) != false || A.FinalStates().Probe(u) != true {
					t.Error()
				}

				// do not visit c two times, or π two times
				if lastLetter == nil {
					lastLetter = x
				} else {
					if lastLetter.IsEqual(x) != false {
						t.Error()
					}
				}

				visitedRStates.Add(s, u)
			}

			transitionCount += 1
		}
	}

	if transitionCount != 3 || visitedRStates.Size() != 4 {
		t.Error()
	}
}

func TestConcatNfa2(t *testing.T) {
	a := nfa.Letter("a")

	A := regex.Concat(regex.Letter("a"), regex.Letter("a")).Nfa()

	if A.Alphabet().Size() != 1 || A.States().Size() < 3 || A.States().Size() > 4 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}
	if A.Alphabet().Probe(a) != true {
		t.Error()
	}

	q0_, _ := A.InitialStates().At(0)
	q0 := q0_.(nfa.State)

	Q0 := A.Transition(q0, a)
	if Q0.Size() < 1 || Q0.Size() > 2 {
		t.Error()
	}

	q1_, _ := Q0.At(0)
	q1 := q1_.(nfa.State)
	if A.Transition(q1, a).Size() == 0 {
		q1_, _ = Q0.At(1)
		q1 = q1_.(nfa.State)
	}

	Q1 := A.Transition(q1, a)
	if Q1.Size() != 1 {
		t.Error()
	}

	if set.Intersect(Q1, A.FinalStates()).Size() != 1 {
		t.Error()
	}
}

func TestLetterNfa(t *testing.T) {

	a := nfa.Letter("a")

	A := regex.Letter("a").Nfa()

	if A.Alphabet().Size() != 1 || A.States().Size() != 2 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}
	if A.Alphabet().Probe(a) != true {
		t.Error()
	}

	//has exactly one transition which is not a loop,
	//and goes from an initial to a final state
	for i := 0; i < A.States().Size(); i += 1 {
		s, _ := A.States().At(i)
		for j := 0; j < A.Alphabet().Size(); j += 1 {
			x, _ := A.Alphabet().At(j)
			S := A.Transition(s.(nfa.State), x.(nfa.Letter))

			if S.Size() == 0 {
				continue
			}
			if S.Size() != 1 {
				t.Error()
			}

			u, _ := S.At(0)

			if s.IsEqual(u) || A.InitialStates().Probe(s) != true || A.FinalStates().Probe(s) != false || A.FinalStates().Probe(u) != true || A.InitialStates().Probe(u) != false {
				t.Error()
			}
		}
	}
}

func TestOrNfa(t *testing.T) {
	a := nfa.Letter("a")
	b := nfa.Letter("π")

	A := regex.Or(regex.Letter("a"), regex.Letter("π")).Nfa()

	if A.Alphabet().Size() != 2 || A.States().Size() != 4 || A.InitialStates().Size() != 2 || A.FinalStates().Size() != 2 {
		t.Error()
	}
	if A.Alphabet().Probe(a) != true || A.Alphabet().Probe(b) != true {
		t.Error()
	}

	//has exactly two transitions which are not loops,
	//both go from an initial to a final state,
	//and use two different letters
	transitionCount := 0
	var lastLetter set.Element
	for i := 0; i < A.States().Size(); i += 1 {
		s, _ := A.States().At(i)
		for j := 0; j < A.Alphabet().Size(); j += 1 {
			x, _ := A.Alphabet().At(j)
			S := A.Transition(s.(nfa.State), x.(nfa.Letter))

			if S.Size() == 0 {
				continue
			}
			if S.Size() != 1 {
				t.Error()
			}

			u, _ := S.At(0)

			if s.IsEqual(u) || A.InitialStates().Probe(s) != true || A.FinalStates().Probe(s) != false || A.FinalStates().Probe(u) != true || A.InitialStates().Probe(u) != false {
				t.Error()
			}

			if transitionCount == 0 {
				lastLetter = x
			} else if transitionCount == 1 {
				if lastLetter.IsEqual(x) {
					t.Error()
				}
			}

			transitionCount += 1
		}
	}

	if transitionCount != 2 {
		t.Error()
	}
}

func TestStarNfa(t *testing.T) {

	a := nfa.Letter("a")
	b := nfa.Letter("π")

	A := regex.Star(regex.Or(regex.Concat(regex.Letter("a"), regex.Letter("a")), regex.Concat(regex.Letter("π"), regex.Letter("π")))).Nfa()

	// the maximal nfa has 9 states and a lot of transitions to useless states
	//
	// minimal nfa:
	//
	// +-- a ---
	// |        \
	// | - a -> o
	// v/
	// □
	// ^\
	// | - π -> o
	// |       /
	// +-- π --

	if A.Alphabet().Size() != 2 || A.States().Size() < 3 || A.States().Size() > 9 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}
	if A.Alphabet().Probe(a) != true || A.Alphabet().Probe(b) != true {
		t.Error()
	}
	if A.InitialStates().IsEqual(A.FinalStates()) != true {
		t.Error()
	}

	q0_, _ := A.InitialStates().At(0)
	q0 := q0_.(nfa.State)

	if A.Transition(q0, a).Size() < 1 || A.Transition(q0, a).Size() > 2 || A.Transition(q0, b).Size() < 1 || A.Transition(q0, b).Size() > 2 {
		t.Error()
	}

	q1_, _ := A.Transition(q0, a).At(0)
	q1 := q1_.(nfa.State)
	if A.Transition(q1, a).Size() == 0 {
		q1_, _ = A.Transition(q0, a).At(1)
		q1 = q1_.(nfa.State)
	}

	q2_, _ := A.Transition(q0, b).At(0)
	q2 := q2_.(nfa.State)
	if A.Transition(q2, b).Size() == 0 {
		q2_, _ = A.Transition(q0, b).At(1)
		q2 = q2_.(nfa.State)
	}

	if A.Transition(q1, a).Size() < 1 || A.Transition(q1, a).Size() > 3 || A.Transition(q2, b).Size() < 1 || A.Transition(q2, b).Size() > 3 {
		t.Error()
	}

	Q1 := A.Transition(q1, a)
	Q2 := A.Transition(q2, b)

	if Q1.Probe(q0) == false || Q2.Probe(q0) == false {
		t.Error()
	}

	var q10, q20 nfa.State

	if q, _ := Q1.At(0); q.IsEqual(q0) {
		q10_, _ := Q1.At(1)
		q10 = q10_.(nfa.State)
	} else {
		q10_, _ := Q1.At(0)
		q10 = q10_.(nfa.State)
	}
	if q, _ := Q2.At(0); q.IsEqual(q0) {
		q20_, _ := Q2.At(1)
		q20 = q20_.(nfa.State)
	} else {
		q20_, _ := Q2.At(0)
		q20 = q20_.(nfa.State)
	}

	if A.Transition(q10, a).Size() != 0 || A.Transition(q10, b).Size() != 0 || A.Transition(q20, a).Size() != 0 || A.Transition(q20, b).Size() != 0 {
		t.Error()
	}
}
