package regex

import (
	"fmt"
	"github.com/hydroo/gomochex/automaton/nfa"
	"github.com/hydroo/gomochex/basic/set"
	"strings"
	"unicode/utf8"
)

type Expression interface {
	String() string
	IsEqual(Expression) bool
	Nfa() nfa.Nfa
}

/*****************************************************************************/

type concatExpression struct {
	l, r Expression
}

func (e concatExpression) String() string {
	return fmt.Sprint("(", e.l, ".", e.r, ")")
}

func (e concatExpression) IsEqual(f_ Expression) bool {
	if f, ok := f_.(concatExpression); ok == true {
		return e.l.IsEqual(f.l) && e.r.IsEqual(f.r)
	} // else {
	return false
	//}
}

func (e concatExpression) Nfa() nfa.Nfa {
	L := e.l.Nfa()
	R := e.r.Nfa()
	A := nfa.NewNfa()

	A.SetAlphabet(set.Join(L.Alphabet(), R.Alphabet()))

	//add R as it is with a 1 prepended to all states
	//add L as it is with a 0 prepended to all states
	for k, T := range []nfa.Nfa{L, R} {
		for i := 0; i < T.States().Size(); i += 1 {
			s, _ := T.States().At(i)
			s_ := nfa.State(fmt.Sprint(k, s))

			A.States().Add(s_)

			if k == 1 && T.FinalStates().Probe(s) == true { //R
				A.FinalStates().Add(s_)
			} else if k == 0 && T.InitialStates().Probe(s) == true { //L
				A.InitialStates().Add(s_)
			}

			for j := 0; j < T.Alphabet().Size(); j += 1 {
				a, _ := T.Alphabet().At(j)
				S := T.Transition(s.(nfa.State), a.(nfa.Letter))

				S_ := set.NewSet()
				for l := 0; l < S.Size(); l += 1 {
					t, _ := S.At(l)
					S_.Add(nfa.State(fmt.Sprint(k, t)))
				}

				A.SetTransition(s_, a.(nfa.Letter), S_)
			}
		}
	}

	// add transititons from before final states in L to all initial states in R
	S_ := set.NewSet()
	for i := 0; i < R.InitialStates().Size(); i += 1 {
		s, _ := R.InitialStates().At(i)
		S_.Add(nfa.State(fmt.Sprint(1, s)))
	}

	for i := 0; i < L.States().Size(); i += 1 {
		s, _ := L.States().At(i)
		s_ := nfa.State(fmt.Sprint(0, s))

		for j := 0; j < L.Alphabet().Size(); j += 1 {
			a, _ := L.Alphabet().At(j)
			S := L.Transition(s.(nfa.State), a.(nfa.Letter))

			if set.Intersect(L.FinalStates(), S).Size() > 0 {
				A.SetTransition(s_, a.(nfa.Letter), set.Join(A.Transition(s_, a.(nfa.Letter)), S_))
			}
		}
	}

	return A
}

type letterExpression struct {
	l string
}

func (e letterExpression) String() string {
	return e.l
}

func (e letterExpression) IsEqual(f_ Expression) bool {
	if f, ok := f_.(letterExpression); ok == true {
		return e == f
	} // else {
	return false
	//}
}

func (e letterExpression) Nfa() nfa.Nfa {
	A := nfa.NewNfa()

	l := nfa.Letter(e.l)
	q0 := nfa.State("0")
	qf := nfa.State("f")

	A.Alphabet().Add(l)
	A.States().Add(q0, qf)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(qf)

	A.SetTransition(q0, l, set.NewSet(qf))

	return A
}

type orExpression struct {
	l, r Expression
}

func (e orExpression) String() string {
	return fmt.Sprint("(", e.l, "+", e.r, ")")
}

func (e orExpression) IsEqual(f_ Expression) bool {
	if f, ok := f_.(orExpression); ok == true {
		return (e.l.IsEqual(f.l) && e.r.IsEqual(f.r)) || (e.l.IsEqual(f.r) && e.r.IsEqual(f.l))
	} // else {
	return false
	//}
}

func (e orExpression) Nfa() nfa.Nfa {
	L := e.l.Nfa()
	R := e.r.Nfa()
	A := nfa.NewNfa()

	A.SetAlphabet(set.Join(L.Alphabet(), R.Alphabet()))

	for k, T := range []nfa.Nfa{L, R} {

		for i := 0; i < T.States().Size(); i += 1 {
			s, _ := T.States().At(i)
			ss := nfa.State(fmt.Sprint(k, s))
			A.States().Add(ss)

			if T.InitialStates().Probe(s) == true {
				A.InitialStates().Add(ss)
			}
			if T.FinalStates().Probe(s) == true {
				A.FinalStates().Add(ss)
			}

			for j := 0; j < T.Alphabet().Size(); j += 1 {
				a, _ := T.Alphabet().At(j)

				S := T.Transition(s.(nfa.State), a.(nfa.Letter))
				SS := set.NewSet()
				for l := 0; l < S.Size(); l += 1 {
					t, _ := S.At(l)
					SS.Add(nfa.State(fmt.Sprint(k, t)))
				}

				A.SetTransition(ss, a.(nfa.Letter), SS)
			}
		}

	}

	return A
}

type starExpression struct {
	f Expression
}

func (e starExpression) String() string {
	return fmt.Sprint("(", e.f, ")*")
}

func (e starExpression) IsEqual(f_ Expression) bool {
	if f, ok := f_.(starExpression); ok == true {
		return e.f.IsEqual(f.f)
	} // else {
	return false
	//}
}

func (e starExpression) Nfa() nfa.Nfa {

	A := e.f.Nfa()

	var q0 nfa.State
	for i := 0; ; i += 1 {
		q0 = nfa.State(fmt.Sprint(i))
		if A.States().Probe(q0) == false {
			break
		}
	}

	A.States().Add(q0)

	// add transitions from before final states to the new initial state
	for i := 0; i < A.States().Size(); i += 1 {
		for j := 0; j < A.Alphabet().Size(); j += 1 {

			q_, _ := A.States().At(i)
			q := q_.(nfa.State)

			a_, _ := A.Alphabet().At(j)
			a := a_.(nfa.Letter)

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
			q := q_.(nfa.State)

			a_, _ := A.Alphabet().At(j)
			a := a_.(nfa.Letter)

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

/*****************************************************************************/

func Concat(l, r Expression) Expression {
	return concatExpression{l, r}
}

func Letter(l string) Expression {
	return letterExpression{l}
}

func Or(l, r Expression) Expression {
	return orExpression{l, r}
}

func Star(e Expression) Expression {
	return starExpression{e}
}

/*****************************************************************************/

func ExpressionFromString(s string) (Expression, bool) {
	return expressionFromStringRecursively(strings.Replace(s, " ", "", -1))
}

func expressionFromStringRecursively(s string) (Expression, bool) {
	firstRune, _ := utf8.DecodeRune([]byte(s))

	if firstRune == '(' {
		bracketCount := 0
		for i := 1; i < len(s); {
			r, runeSize := utf8.DecodeRune([]byte(s[i:]))
			if r == '(' {
				bracketCount += 1
			} else if r == ')' {
				bracketCount -= 1
			} else if bracketCount == 0 && (r == '.' || r == '+') {
				subL, okL := expressionFromStringRecursively(s[1:i])
				subR, okR := expressionFromStringRecursively(s[i+runeSize : len(s)-1])
				if r == '.' {
					return Concat(subL, subR), okL && okR
				} else { // '+'
					return Or(subL, subR), okL && okR
				}
			} else if bracketCount == -1 && r == '*' { // star
				sub, ok := expressionFromStringRecursively(s[1 : i-1])
				return Star(sub), ok
			} else if bracketCount == -1 && r != '*' {
				return nil, false
			}

			i += runeSize
		}

	} else { //letter
		for i := 0; i < len(s); {
			r, runeSize := utf8.DecodeRune([]byte(s[i:]))
			if r == '(' || r == ')' || r == '.' || r == '+' || r == '*' {
				return nil, false
			}
			i += runeSize
		}
		return Letter(s), true
	}

	return nil, false
}
