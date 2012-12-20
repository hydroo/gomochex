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

type concat struct {
	l, r Expression
}

func (e concat) String() string {
	return fmt.Sprint("(", e.l, ".", e.r, ")")
}

func (e concat) Nfa() nfa.Nfa {
	//TODO
	ret := nfa.NewNfa()
	return ret
}

type letter struct {
	l string
}

func (e letter) String() string {
	return fmt.Sprint(e.l)
}

func (e letter) Nfa() nfa.Nfa {
	ret := nfa.NewNfa()

	l := nfa.StringLetter(e.l)
	q0 := nfa.StringState("0")
	qf := nfa.StringState("f")

	ret.Alphabet().Add(l)
	ret.States().Add(q0, qf)
	ret.InitialStates().Add(q0)
	ret.FinalStates().Add(qf)

	trans := func(s nfa.State, m nfa.Letter) nfa.StateSet {
		S := set.NewSet()
		if s.IsEqual(q0) && m.IsEqual(l) {
			S.Add(qf)
		}
		return S
	}

	ret.SetTransitionFunction(trans)

	return ret
}

type or struct {
	l, r Expression
}

func (e or) String() string {
	return fmt.Sprint("(", e.l, "+", e.r, ")")
}

func (e or) Nfa() nfa.Nfa {
	//TODO
	ret := nfa.NewNfa()
	return ret
}

type star struct {
	f Expression
}

func (e star) String() string {
	return fmt.Sprint("(", e.f, ")*")
}

func (e star) Nfa() nfa.Nfa {
	//TODO
	ret := nfa.NewNfa()
	return ret
}

/*****************************************************************************/

func Concat(l, r Expression) Expression {
	return &concat{l, r}
}

func Letter(l string) Expression {
	return &letter{l}
}

func Or(l, r Expression) Expression {
	return &or{l, r}
}

func Star(e Expression) Expression {
	return &star{e}
}

/*****************************************************************************/

func ExpressionFromString(s string) (Expression, bool) {
	s = strings.Replace(s, " ", "", -1)
	return expressionFromStringRecursively(s)
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
