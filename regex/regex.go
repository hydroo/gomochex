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
	Nfa() nfa.Nfa
}

/*****************************************************************************/

type concatExpression struct {
	l, r Expression
}

func (e concatExpression) String() string {
	return fmt.Sprint("(", e.l, ".", e.r, ")")
}

func (e concatExpression) Nfa() nfa.Nfa {
	//TODO
	ret := nfa.NewNfa()
	return ret
}

type letterExpression struct {
	l string
}

func (e letterExpression) String() string {
	return e.l
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

func (e orExpression) Nfa() nfa.Nfa {
	L := e.l.Nfa()
	R := e.r.Nfa()
	A := nfa.NewNfa()

	A.SetAlphabet(set.Join(L.Alphabet(), R.Alphabet()))

	for k, T := range []nfa.Nfa{L, R} {

		for i := 0; i < T.States().Size(); i += 1 {
			s, _ := T.States().At(i)
			ss := nfa.State(fmt.Sprint(k,s))
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
					SS.Add(nfa.State(fmt.Sprint(k,t)))
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

func (e starExpression) Nfa() nfa.Nfa {
	//TODO
	ret := nfa.NewNfa()
	return ret
}

/*****************************************************************************/

func Concat(l, r Expression) Expression {
	return &concatExpression{l, r}
}

func Letter(l string) Expression {
	return &letterExpression{l}
}

func Or(l, r Expression) Expression {
	return &orExpression{l, r}
}

func Star(e Expression) Expression {
	return &starExpression{e}
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
