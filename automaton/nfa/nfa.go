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
	json.Marshaler
	json.Unmarshaler

	Alphabet() Alphabet
	SetAlphabet(Alphabet)

	InitialStates() StateSet
	SetInitialStates(StateSet)

	FinalStates() StateSet
	SetFinalStates(StateSet)

	States() StateSet
	SetStates(StateSet)

	String() string

	Transition(State, Letter) StateSet
	SetTransition(State, Letter, StateSet)
	SetTransitionFunction(func(State, Letter) StateSet)

	Copy() Nfa
}

func NewNfa() Nfa {
	return &simpleNfa{set.NewSet(), set.NewSet(), set.NewSet(), make(map[State]map[Letter]StateSet), set.NewSet()}
}

/*****************************************************************************/

type simpleNfa struct {
	states        StateSet
	alphabet      set.Set
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
		Alphabet      set.Set
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
	B.SetAlphabet(A.Alphabet().Copy().(set.Set))
	B.SetInitialStates(A.InitialStates().Copy().(StateSet))
	B.SetFinalStates(A.FinalStates().Copy().(StateSet))

	for k, v := range A.transitions {
		for l, w := range v {
			B.SetTransition(k, l, w.Copy().(StateSet))
		}
	}

	return B
}
