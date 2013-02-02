package nfa

import (
	"encoding/json"
	"fmt"
	"github.com/hydroo/gomochex/basic/set"
)

type State string

func (s State) IsEqual(e set.Element) bool {
	f, ok := e.(State)
	if ok == false {
		return false
	} //else {
	return s == f
	//}
}

type Letter string

func (s Letter) IsEqual(e set.Element) bool {
	f, ok := e.(Letter)
	if ok == false {
		return false
	} //else {
	return s == f
	//}
}

type Alphabet set.Set

type StateSet set.Set

type Nfa interface {
	Alphabet() Alphabet
	SetAlphabet(Alphabet)

	InitialStates() StateSet
	SetInitialStates(StateSet)

	FinalStates() StateSet
	SetFinalStates(StateSet)

	States() StateSet
	SetStates(StateSet)

	Transition(State, Letter) StateSet
	SetTransition(State, Letter, StateSet)
	SetTransitionFunction(func(State, Letter) StateSet)

	Copy() Nfa


	String() string
	json.Marshaler
	json.Unmarshaler
}

func NewNfa() Nfa {
	return &simpleNfa{set.NewSet(), set.NewSet(), set.NewSet(), make(map[State]map[Letter]StateSet), set.NewSet()}
}

func Concat(A, B Nfa) Nfa {
	C := NewNfa()

	C.SetAlphabet(set.Join(A.Alphabet(), B.Alphabet()))

	//add B as it is with a 1 prepended to all states
	//add A as it is with a 0 prepended to all states
	for k, T := range []Nfa{A, B} {
		for i := 0; i < T.States().Size(); i += 1 {
			s, _ := T.States().At(i)
			s_ := State(fmt.Sprint(k, s))

			C.States().Add(s_)

			if k == 1 && T.FinalStates().Probe(s) == true { //B
				C.FinalStates().Add(s_)
			} else if k == 0 && T.InitialStates().Probe(s) == true { //A
				C.InitialStates().Add(s_)
			}

			for j := 0; j < T.Alphabet().Size(); j += 1 {
				a, _ := T.Alphabet().At(j)
				S := T.Transition(s.(State), a.(Letter))

				S_ := set.NewSet()
				for l := 0; l < S.Size(); l += 1 {
					t, _ := S.At(l)
					S_.Add(State(fmt.Sprint(k, t)))
				}

				C.SetTransition(s_, a.(Letter), S_)
			}
		}
	}

	// add transititons from before final states in A to all initial states in B
	S_ := set.NewSet()
	for i := 0; i < B.InitialStates().Size(); i += 1 {
		s, _ := B.InitialStates().At(i)
		S_.Add(State(fmt.Sprint(1, s)))
	}

	for i := 0; i < A.States().Size(); i += 1 {
		s, _ := A.States().At(i)
		s_ := State(fmt.Sprint(0, s))

		for j := 0; j < A.Alphabet().Size(); j += 1 {
			a, _ := A.Alphabet().At(j)
			S := A.Transition(s.(State), a.(Letter))

			if set.Intersect(A.FinalStates(), S).Size() > 0 {
				C.SetTransition(s_, a.(Letter), set.Join(A.Transition(s_, a.(Letter)), S_))
			}
		}
	}

	return C
}

func KleeneStar(A Nfa) Nfa {

	var q0 State
	for i := 0; ; i += 1 {
		q0 = State(fmt.Sprint(i))
		if A.States().Probe(q0) == false {
			break
		}
	}

	A.States().Add(q0)

	// add transitions from before final states to the new initial state
	for i := 0; i < A.States().Size(); i += 1 {
		for j := 0; j < A.Alphabet().Size(); j += 1 {

			q_, _ := A.States().At(i)
			q := q_.(State)

			a_, _ := A.Alphabet().At(j)
			a := a_.(Letter)

			Q := A.Transition(q, a)

			if set.Intersect(Q, A.FinalStates()).Size() != 0 {
				Q.Add(q0)
				A.SetTransition(q, a, Q)
			}
		}
	}

	// add transitions from the new initial state
	// to the destinations of the old initial states
	for i := 0; i < A.InitialStates().Size(); i += 1 {
		for j := 0; j < A.Alphabet().Size(); j += 1 {

			q_, _ := A.InitialStates().At(i)
			q := q_.(State)

			a_, _ := A.Alphabet().At(j)
			a := a_.(Letter)

			Q := A.Transition(q, a)
			Q0 := A.Transition(q0, a)

			A.SetTransition(q0, a, set.Join(Q, Q0))
		}
	}

	A.InitialStates().Clear()
	A.InitialStates().Add(q0)
	A.FinalStates().Clear()
	A.FinalStates().Add(q0)

	return A
}

func OneLetter(a Letter) Nfa {
	A := NewNfa()
	q0 := State("0")
	qf := State("f")
	A.Alphabet().Add(a)
	A.States().Add(q0, qf)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(qf)
	A.SetTransition(q0, a, set.NewSet(qf))
	return A
}

func Union(A, B Nfa) Nfa {
	C := NewNfa()

	C.SetAlphabet(set.Join(A.Alphabet(), B.Alphabet()))

	for k, T := range []Nfa{A, B} {
		for i := 0; i < T.States().Size(); i += 1 {
			s, _ := T.States().At(i)
			ss := State(fmt.Sprint(k, s))
			C.States().Add(ss)

			if T.InitialStates().Probe(s) == true {
				C.InitialStates().Add(ss)
			}
			if T.FinalStates().Probe(s) == true {
				C.FinalStates().Add(ss)
			}

			for j := 0; j < T.Alphabet().Size(); j += 1 {
				a, _ := T.Alphabet().At(j)

				S := T.Transition(s.(State), a.(Letter))
				SS := set.NewSet()
				for l := 0; l < S.Size(); l += 1 {
					t, _ := S.At(l)
					SS.Add(State(fmt.Sprint(k, t)))
				}

				C.SetTransition(ss, a.(Letter), SS)
			}
		}
	}

	return C
}

/*****************************************************************************/

type simpleNfa struct {
	states        StateSet
	alphabet      Alphabet
	initialStates StateSet
	transitions   map[State]map[Letter]StateSet
	finalStates   StateSet
}

func (A simpleNfa) Alphabet() Alphabet {
	return A.alphabet
}

func (A *simpleNfa) SetAlphabet(sigma Alphabet) {
	A.alphabet = sigma
}

func (A simpleNfa) InitialStates() StateSet {
	return A.initialStates
}

func (A *simpleNfa) SetInitialStates(S StateSet) {
	A.initialStates = S
}

func (A simpleNfa) FinalStates() StateSet {
	return A.finalStates
}

func (A *simpleNfa) SetFinalStates(F StateSet) {
	A.finalStates = F
}

func (A simpleNfa) States() StateSet {
	return A.states
}

func (A *simpleNfa) SetStates(S StateSet) {
	A.states = S
}

func (A simpleNfa) String() string {
	ret := ""
	ret += fmt.Sprintln("states:", A.States())
	ret += fmt.Sprintln("alphabet:", A.Alphabet())
	ret += fmt.Sprintln("initial states:", A.InitialStates())
	ret += fmt.Sprintln("final states:", A.FinalStates())
	ret += fmt.Sprintln("transitions:")
	for i := 0; i < A.States().Size(); i += 1 {
		for j := 0; j < A.Alphabet().Size(); j += 1 {
			s, _ := A.States().At(i)
			a, _ := A.Alphabet().At(j)
			if S := A.Transition(s.(State), a.(Letter)); S.Size() > 0 {
				ret += fmt.Sprintln(" ", s, "--", a, "-->", S)
			}
		}
	}

	return ret
}

func (A simpleNfa) MarshalJSON() ([]byte, error) {
	type simpleNfaWithExportedFields struct {
		States        StateSet
		Alphabet      Alphabet
		InitialStates StateSet
		Transitions   map[State]map[Letter]StateSet
		FinalStates   StateSet
	}
	return json.Marshal(simpleNfaWithExportedFields{A.states, A.alphabet, A.initialStates, A.transitions, A.finalStates})
}

func (A *simpleNfa) UnmarshalJSON(b []byte) error {
	type simpleNfaForUnmarshaling struct {
		States        []State
		Alphabet      []Letter
		InitialStates []State
		Transitions   map[string]map[string][]string //using Letter or State won't work here
		FinalStates   []State
	}

	var B simpleNfaForUnmarshaling
	if err := json.Unmarshal(b, &B); err != nil {
		return err
	}

	states := set.NewSet()
	alphabet := set.NewSet()
	initialStates := set.NewSet()
	transitions := make(map[State]map[Letter]StateSet)
	finalStates := set.NewSet()

	for _, s := range B.States {
		states.Add(s)
	}
	for _, l := range B.Alphabet {
		alphabet.Add(l)
	}
	for _, s := range B.InitialStates {
		initialStates.Add(s)
	}
	for _, s := range B.FinalStates {
		finalStates.Add(s)
	}
	for k, v := range B.Transitions {
		for l, w := range v {
			if _, ok := transitions[State(k)]; ok != true {
				transitions[State(k)] = make(map[Letter]StateSet)
			}
			S := set.NewSet()
			for _, s := range w {
				S.Add(State(s))
			}
			transitions[State(k)][Letter(l)] = S
		}
	}

	*A = simpleNfa{states, alphabet, initialStates, transitions, finalStates}

	return nil
}

func (A simpleNfa) Transition(s State, l Letter) StateSet {

	var m map[Letter]StateSet
	if _, ok := A.transitions[s]; ok == false {
		return set.NewSet()
	} else {
		m = A.transitions[s]
	}

	var S StateSet
	if _, ok := m[l]; ok == false {
		S = set.NewSet()
	} else {
		S = m[l]
	}

	return S
}

func (A *simpleNfa) SetTransition(s State, l Letter, S StateSet) {

	var m map[Letter]StateSet
	if _, ok := A.transitions[s]; ok == false {
		m = make(map[Letter]StateSet)
		A.transitions[s] = m
	} else {
		m = A.transitions[s]
	}

	if S.Size() == 0 {
		delete(m, l)
	} else {
		m[l] = S
	}
}

func (A *simpleNfa) SetTransitionFunction(delta func(State, Letter) StateSet) {

	A.transitions = make(map[State]map[Letter]StateSet)

	for i := 0; i < A.States().Size(); i += 1 {
		for j := 0; j < A.Alphabet().Size(); j += 1 {

			q, _ := A.States().At(i)
			a, _ := A.Alphabet().At(j)

			Q := delta(q.(State), a.(Letter))

			if Q.Size() > 0 {
				A.SetTransition(q.(State), a.(Letter), Q)
			}
		}
	}
}

func (A simpleNfa) Copy() Nfa {
	B := NewNfa()
	B.SetStates(A.States().Copy().(StateSet))
	B.SetAlphabet(A.Alphabet().Copy().(Alphabet))
	B.SetInitialStates(A.InitialStates().Copy().(StateSet))
	B.SetFinalStates(A.FinalStates().Copy().(StateSet))

	for k, v := range A.transitions {
		for l, w := range v {
			B.SetTransition(k, l, w.Copy().(StateSet))
		}
	}

	return B
}

// Q will be the new initial states
func (A simpleNfa) inducedNfa(q State) Nfa {
	B := NewNfa().(*simpleNfa)
	B.initialStates = set.NewSet(q)
	B.alphabet = A.Alphabet()

	var recurse func(State)
	recurse = func(q State) {
		if B.States().Probe(q) == true {
			return
		}

		B.States().Add(q)

		if B.transitions[q] == nil && len(A.transitions[q]) > 0 {
			B.transitions[q] = make(map[Letter]StateSet)
		}

		for l, w := range A.transitions[q] {
			B.transitions[q][l] = w.Copy().(StateSet)

			for i := 0; i < w.Size(); i += 1 {
				r, _ := w.At(i)
				recurse(r.(State))
			}
		}
	}

	recurse(q)

	B.finalStates = set.Intersect(A.FinalStates(), B.States())

	return B
}

func (A simpleNfa) reachableStates() StateSet {
	reached := set.NewSet()

	var recurse func(State)
	recurse = func(q State) {
		if reached.Probe(q) == true {
			return
		}

		reached.Add(q)

		for i := 0; i < A.Alphabet().Size(); i += 1 {
			a, _ := A.Alphabet().At(i)
			Q := A.Transition(q, a.(Letter))

			for j := 0; j < Q.Size(); j += 1 {
				q_, _ := Q.At(j)
				recurse(q_.(State))
			}
		}
	}

	for i := 0; i < A.InitialStates().Size(); i += 1 {
		q, _ := A.InitialStates().At(i)
		recurse(q.(State))
	}

	return reached
}

// states and transitions from which no final state is reachable
// states and transitions that are unreachable
func (A simpleNfa) removeUselessParts() Nfa {
	reachableStates := A.reachableStates()

	nonEmptyLanguageStates := set.NewSet()

	var recurse func(State)
	recurse = func(u State) {
		if nonEmptyLanguageStates.Probe(u) == true {
			return
		}
		nonEmptyLanguageStates.Add(u)
		for k, v := range A.transitions { //map[State]map[Letter]StateSet
			for _, w := range v {
				if set.Intersect(set.NewSet(u), w).Size() > 0 {
					recurse(k)
				}
			}
		}
	}

	//backward search for all states that can reach a reachable final state
	reachableFinalStates := set.Intersect(reachableStates, A.FinalStates())
	for i := 0; i < reachableFinalStates.Size(); i += 1 {
		q, _ := reachableFinalStates.At(i)
		recurse(q.(State))
	}

	usefulStates := set.Intersect(reachableFinalStates, nonEmptyLanguageStates)

	B := NewNfa().(*simpleNfa)

	B.states = usefulStates
	B.initialStates = set.Intersect(A.InitialStates(), usefulStates)
	B.finalStates = set.Intersect(A.InitialStates(), usefulStates)

	for k, v := range A.transitions {
		for l, w := range v {
			usefulGoalStates := set.Intersect(usefulStates, w)
			if usefulStates.Probe(k) == true && usefulGoalStates.Size() > 0 {
				B.alphabet.Add(l)
				B.SetTransition(k, l, usefulGoalStates)
			}
		}
	}

	return B
}
