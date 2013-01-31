package nfa

import (
	"bytes"
	"encoding/json"
	//"fmt"
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

	a := Letter("a")
	b := Letter("b")
	q0 := State("1")
	q1 := State("2")
	q2 := State("3")

	A := NewNfa()

	A.Alphabet().Add(a, b)
	A.States().Add(q0, q1, q2)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(q2)

	trans := func(s State, l Letter) StateSet {
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

func TestJson(t *testing.T) {
	a := Letter("a")
	q0 := State("0")
	q1 := State("1")
	q2 := State("2")
	q3 := State("3")

	// -> □ -> o
	//
	//    o -> □
	A := NewNfa()
	A.Alphabet().Add(a)
	A.States().Add(q0, q1, q2, q3)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(q0, q3)
	A.SetTransition(q0, a, set.NewSet(q1))
	A.SetTransition(q2, a, set.NewSet(q3))

	s1, err1 := json.Marshal(A)
	if err1 != nil {
		t.Error(err1)
	}

	B := NewNfa()
	err2 := json.Unmarshal(s1, &B)
	if err2 != nil {
		t.Error(err2)
	}

	s3, err3 := json.Marshal(B)
	if err1 != nil {
		t.Error(err3)
	}

	if bytes.Compare(s1, s3) != 0 {
		t.Error()
	}
}
