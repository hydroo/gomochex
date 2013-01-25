package nfa_test

import (
	//"fmt"
	"github.com/hydroo/gomochex/automaton/nfa"
	"github.com/hydroo/gomochex/basic/set"
	"testing"
)

//                           +-+
// --> 1 -- a --> 2 -- b --> |3|
//     ^                     +-+
//     |                      |
//     +--------- a ----------+
//
func TestSimpleNfa(t *testing.T) {

	a := nfa.Letter("a")
	b := nfa.Letter("b")
	q0 := nfa.State("1")
	q1 := nfa.State("2")
	q2 := nfa.State("3")

	A := nfa.NewNfa()

	A.Alphabet().Add(a, b)
	A.States().Add(q0, q1, q2)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(q2)

	trans := func(s nfa.State, l nfa.Letter) nfa.StateSet {
		S := set.NewSet()
		if s.IsEqual(q0) && l.IsEqual(a) {
			S.Add(q1)
		} else if s.IsEqual(q1) && l.IsEqual(b) {
			S.Add(q2)
		} else if s.IsEqual(q2) && l.IsEqual(a) {
			S.Add(q0)
		}
		return S
	}

	A.SetTransitionFunction(trans)

	if A.Alphabet().Size() != 2 || A.States().Size() != 3 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}

	if A.Transition(q0, a).IsEqual(set.NewSet(q1)) != true || A.Transition(q1, b).IsEqual(set.NewSet(q2)) != true || A.Transition(q2, a).IsEqual(set.NewSet(q0)) != true || A.Transition(q0, b).IsEqual(set.NewSet()) != true {
		t.Error()
	}
}
