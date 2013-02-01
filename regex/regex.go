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

func Concat(l, r Expression) Expression {
	return concatExpression{l, r}
}

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
	return nfa.Concat(e.l.Nfa(), e.r.Nfa())
}

func Letter(l string) Expression {
	return letterExpression{l}
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
	return nfa.OneLetter(nfa.Letter(e.l))
}

func Or(l, r Expression) Expression {
	return orExpression{l, r}
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
	return nfa.Union(e.l.Nfa(), e.r.Nfa())
}

func Star(e Expression) Expression {
	return starExpression{e}
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
	return nfa.KleeneStar(e.f.Nfa())
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
