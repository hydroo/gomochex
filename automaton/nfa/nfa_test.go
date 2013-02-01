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

func TestConcat(t *testing.T) {
	a := Letter("a")
	b := Letter("π")
	c := Letter("c")
	A := Concat(OneLetter(a), Union(OneLetter(b), OneLetter(c)))

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
			S := A.Transition(s.(State), x.(Letter))

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
	a := Letter("a")
	A := Concat(OneLetter(a), OneLetter(a))

	if A.Alphabet().Size() != 1 || A.States().Size() < 3 || A.States().Size() > 4 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}
	if A.Alphabet().Probe(a) != true {
		t.Error()
	}

	q0_, _ := A.InitialStates().At(0)
	q0 := q0_.(State)

	Q0 := A.Transition(q0, a)
	if Q0.Size() < 1 || Q0.Size() > 2 {
		t.Error()
	}

	q1_, _ := Q0.At(0)
	q1 := q1_.(State)
	if A.Transition(q1, a).Size() == 0 {
		q1_, _ = Q0.At(1)
		q1 = q1_.(State)
	}

	Q1 := A.Transition(q1, a)
	if Q1.Size() != 1 {
		t.Error()
	}

	if set.Intersect(Q1, A.FinalStates()).Size() != 1 {
		t.Error()
	}
}

func TestKleeneStar(t *testing.T) {
	a := Letter("a")
	b := Letter("π")
	A := KleeneStar(Union(Concat(OneLetter(a), OneLetter(a)), Concat(OneLetter(b), OneLetter(b))))

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
	q0 := q0_.(State)

	if A.Transition(q0, a).Size() < 1 || A.Transition(q0, a).Size() > 2 || A.Transition(q0, b).Size() < 1 || A.Transition(q0, b).Size() > 2 {
		t.Error()
	}

	q1_, _ := A.Transition(q0, a).At(0)
	q1 := q1_.(State)
	if A.Transition(q1, a).Size() == 0 {
		q1_, _ = A.Transition(q0, a).At(1)
		q1 = q1_.(State)
	}

	q2_, _ := A.Transition(q0, b).At(0)
	q2 := q2_.(State)
	if A.Transition(q2, b).Size() == 0 {
		q2_, _ = A.Transition(q0, b).At(1)
		q2 = q2_.(State)
	}

	if A.Transition(q1, a).Size() < 1 || A.Transition(q1, a).Size() > 3 || A.Transition(q2, b).Size() < 1 || A.Transition(q2, b).Size() > 3 {
		t.Error()
	}

	Q1 := A.Transition(q1, a)
	Q2 := A.Transition(q2, b)

	if Q1.Probe(q0) == false || Q2.Probe(q0) == false {
		t.Error()
	}

	var q10, q20 State

	if q, _ := Q1.At(0); q.IsEqual(q0) {
		q10_, _ := Q1.At(1)
		q10 = q10_.(State)
	} else {
		q10_, _ := Q1.At(0)
		q10 = q10_.(State)
	}
	if q, _ := Q2.At(0); q.IsEqual(q0) {
		q20_, _ := Q2.At(1)
		q20 = q20_.(State)
	} else {
		q20_, _ := Q2.At(0)
		q20 = q20_.(State)
	}

	if A.Transition(q10, a).Size() != 0 || A.Transition(q10, b).Size() != 0 || A.Transition(q20, a).Size() != 0 || A.Transition(q20, b).Size() != 0 {
		t.Error()
	}
}

func TestOneLetter(t *testing.T) {
	a := Letter("π")
	A := OneLetter(a)

	if A.Alphabet().Size() != 1 || A.Alphabet().Probe(a) != true || A.States().Size() != 2 || A.InitialStates().Size() != 1 || A.FinalStates().Size() != 1 {
		t.Error()
	}

	//has exactly one transition which is not a loop,
	//and goes from an initial to a final state
	for i := 0; i < A.States().Size(); i += 1 {
		s, _ := A.States().At(i)
		for j := 0; j < A.Alphabet().Size(); j += 1 {
			x, _ := A.Alphabet().At(j)
			S := A.Transition(s.(State), x.(Letter))

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

func TestUnion(t *testing.T) {
	a := Letter("a")
	b := Letter("π")
	A := Union(OneLetter(a), OneLetter(b))

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
			S := A.Transition(s.(State), x.(Letter))

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

