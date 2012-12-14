package ltl

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Formula interface {
	String() string
}

// atomic proposition
type aPFormula struct {
	a string
}

func (n aPFormula) String() string {
	return n.a
}

type notFormula struct {
	phi Formula
}

func (n notFormula) String() string {
	return fmt.Sprint("¬(", n.phi, ")")
}

type andFormula struct {
	phi, psi Formula
}

func (n andFormula) String() string {
	return fmt.Sprint("(", n.phi, "∧" ,n.psi, ")")
}

// atomic proposition
func AP(a string) Formula {
	return aPFormula{a}
}

func Not(phi Formula) Formula {
	return notFormula{phi}
}

func And(phi, psi Formula) Formula {
	return andFormula{phi, psi}
}

func FormulaFromString(phi string) (Formula, bool) {
	phi = strings.Replace(phi, " ", "", -1)
	return formulaFromStringRecursively(phi)
}

func formulaFromStringRecursively(phi string) (Formula, bool) {

	firstRune, firstRuneSize := utf8.DecodeRune([]byte(phi))

	switch {
	case len(phi) == 0 : // error
		return nil, false
	case firstRune == '¬' : // not
		if len(phi)+1 <= firstRuneSize || phi[firstRuneSize] != '(' {
			return nil, false
		}

		phi, ok := formulaFromStringRecursively(phi[3:len(phi)-1])
		return Not(phi), ok
	case firstRune == '(' : // and
		bracketCount := 0
		for i := 1; i < len(phi); {
			b, runeSize := utf8.DecodeRune([]byte(phi[i:]))
			if b == '(' {
				bracketCount += 1
			} else if b == ')' {
				bracketCount -= 1
			} else if bracketCount == 0 && b == '∧' {
				subPhi, okPhi := formulaFromStringRecursively(phi[1:i])
				subPsi, okPsi := formulaFromStringRecursively(phi[i+3:len(phi)-1])
				return And(subPhi, subPsi), okPhi && okPsi
			} else if bracketCount == 0 && phi[i] == ')' {
				return nil, false
			}

			i += runeSize
		}

		return nil, false // too many opening brackets, or no ∧ was found
	case firstRune != '¬' && firstRune != '(' : //ap
		for i := 0; i < len(phi); {
			b, runeSize := utf8.DecodeRune([]byte(phi[i:]))
			if b == '¬' || b == '(' || b == ')' || b == '∧' {
				return nil, false
			}
			i += runeSize
		}

		return AP(phi), true
	}

	return nil, false // error
}

