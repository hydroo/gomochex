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
	s1 := nfa.State("1")
	s2 := nfa.State("2")
	s3 := nfa.State("3")

	A := nfa.NewNfa()

	A.Alphabet().Add(a, b)
	A.States().Add(s1, s2, s3)
	A.InitialStates().Add(s1)
	A.FinalStates().Add(s3)

	trans := func(s nfa.State, l nfa.Letter) nfa.StateSet {
		S := set.NewSet()
		if s.IsEqual(s1) && l.IsEqual(a) {
			S.Add(s2)
		} else if s.IsEqual(s2) && l.IsEqual(b) {
			S.Add(s3)
		} else if s.IsEqual(s3) && l.IsEqual(a) {
			S.Add(s1)
		}
		return S
	}

	A.SetTransitionFunction(trans)

	if A.Alphabet().Size() != 2 || A.States().Size() != 3 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}

	if A.Transition(s1, a).IsEqual(set.NewSet(s2)) != true || A.Transition(s2, b).IsEqual(set.NewSet(s3)) != true || A.Transition(s3, a).IsEqual(set.NewSet(s1)) != true || A.Transition(s1, b).IsEqual(set.NewSet()) != true {
		t.Error()
	}
}
