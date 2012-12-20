package nfa

import (
	"fmt"
	"github.com/hydroo/gomochex/basic/set"
)

type State interface {
	set.Element
}

type Letter interface {
	set.Element
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

	String() string

	Transition(State, Letter) StateSet
	SetTransitionFunction(func(State, Letter) StateSet)
}

func NewNfa() Nfa {
	return &simpleNfa{set.NewSet(),
		set.NewSet(),
		set.NewSet(),
		func(State, Letter) StateSet { return set.NewSet() },
		set.NewSet()}
}

/*****************************************************************************/

type simpleNfa struct {
	states        set.Set
	alphabet      set.Set
	initialStates set.Set
	transition    func(State, Letter) StateSet
	finalStates   set.Set
}

func (A simpleNfa) Alphabet() Alphabet {
	return A.alphabet
}

func (A simpleNfa) SetAlphabet(sigma Alphabet) {
	A.alphabet = sigma
}

func (A simpleNfa) InitialStates() StateSet {
	return A.initialStates
}

func (A simpleNfa) SetInitialStates(S StateSet) {
	A.initialStates = S
}

func (A simpleNfa) FinalStates() StateSet {
	return A.finalStates
}

func (A simpleNfa) SetFinalStates(F StateSet) {
	A.finalStates = F
}

func (A simpleNfa) States() StateSet {
	return A.states
}

func (A simpleNfa) SetStates(S StateSet) {
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
			if next := A.transition(s, a); next.Size() > 0 {
				ret += fmt.Sprintln(" ", s, "--", a, "-->", next)
			}
		}
	}

	return ret
}

func (A simpleNfa) Transition(s State, l Letter) StateSet {
	return A.transition(s, l)
}

func (A *simpleNfa) SetTransitionFunction(delta func(State, Letter) StateSet) {
	A.transition = delta
}

/*****************************************************************************/

type StringState string

func (s StringState) IsEqual(e set.Element) bool {
	f, ok := e.(StringState)
	if ok == false {
		return false
	} //else {
	return s == f
	//}
}

type StringLetter string

func (s StringLetter) IsEqual(e set.Element) bool {
	f, ok := e.(StringLetter)
	if ok == false {
		return false
	} //else {
	return s == f
	//}
}
