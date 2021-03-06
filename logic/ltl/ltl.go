package ltl

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Formula interface {
	String() string
	IsEqual(Formula) bool
}

/*****************************************************************************/

func Always(phi Formula) Formula {
	return alwaysFormula{phi}
}

type alwaysFormula struct {
	phi Formula
}

func (n alwaysFormula) String() string {
	return fmt.Sprint("□(", n.phi, ")")
}

func (e alwaysFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(alwaysFormula); ok == true {
		return e.phi.IsEqual(f.phi)
	} // else {
	return false
	//}
}

func And(phi, psi Formula) Formula {
	return andFormula{phi, psi}
}

type andFormula struct {
	phi, psi Formula
}

func (n andFormula) String() string {
	return fmt.Sprint("(", n.phi, "∧", n.psi, ")")
}

func (e andFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(andFormula); ok == true {
		return (e.phi.IsEqual(f.phi) && e.psi.IsEqual(f.psi)) || (e.phi.IsEqual(f.psi) && e.psi.IsEqual(f.phi))
	} // else {
	return false
	//}
}

// atomic proposition
func Ap(a string) Formula {
	return aPFormula{a}
}

type aPFormula struct {
	a string
}

func (n aPFormula) String() string {
	return n.a
}

func (e aPFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(aPFormula); ok == true {
		return e.a == f.a
	} // else {
	return false
	//}
}

func Eventually(phi Formula) Formula {
	return eventuallyFormula{phi}
}

type eventuallyFormula struct {
	phi Formula
}

func (n eventuallyFormula) String() string {
	return fmt.Sprint("◇(", n.phi, ")")
}

func (e eventuallyFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(eventuallyFormula); ok == true {
		return e.phi.IsEqual(f.phi)
	} // else {
	return false
	//}
}

func False() Formula {
	return falseFormula{}
}

type falseFormula struct {
}

func (n falseFormula) String() string {
	return "false"
}

func (e falseFormula) IsEqual(f_ Formula) bool {
	if _, ok := f_.(falseFormula); ok == true {
		return true
	} // else {
	return false
	//}
}

func Next(phi Formula) Formula {
	return nextFormula{phi}
}

type nextFormula struct {
	phi Formula
}

func (n nextFormula) String() string {
	return fmt.Sprint("○(", n.phi, ")")
}

func (e nextFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(nextFormula); ok == true {
		return e.phi.IsEqual(f.phi)
	} // else {
	return false
	//}
}

func Not(phi Formula) Formula {
	return notFormula{phi}
}

type notFormula struct {
	phi Formula
}

func (n notFormula) String() string {
	return fmt.Sprint("¬(", n.phi, ")")
}

func (e notFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(notFormula); ok == true {
		return e.phi.IsEqual(f.phi)
	} // else {
	return false
	//}
}

func Or(phi, psi Formula) Formula {
	return orFormula{phi, psi}
}

type orFormula struct {
	phi, psi Formula
}

func (n orFormula) String() string {
	return fmt.Sprint("(", n.phi, "∨", n.psi, ")")
}

func (e orFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(orFormula); ok == true {
		return (e.phi.IsEqual(f.phi) && e.psi.IsEqual(f.psi)) || (e.phi.IsEqual(f.psi) && e.psi.IsEqual(f.phi))
	} // else {
	return false
	//}
}

func True() Formula {
	return trueFormula{}
}

type trueFormula struct {
}

func (n trueFormula) String() string {
	return "true"
}

func (e trueFormula) IsEqual(f_ Formula) bool {
	if _, ok := f_.(trueFormula); ok == true {
		return true
	} // else {
	return false
	//}
}

func Until(phi, psi Formula) Formula {
	return untilFormula{phi, psi}
}

type untilFormula struct {
	phi, psi Formula
}

func (n untilFormula) String() string {
	return fmt.Sprint("((", n.phi, ")U(", n.psi, "))")
}

func (e untilFormula) IsEqual(f_ Formula) bool {
	if f, ok := f_.(untilFormula); ok == true {
		return (e.phi.IsEqual(f.phi) && e.psi.IsEqual(f.psi))
	} // else {
	return false
	//}
}

/*****************************************************************************/

func FormulaFromString(phi string) (Formula, bool) {
	return formulaFromStringRecursively(strings.Replace(phi, " ", "", -1))
}

func formulaFromStringRecursively(phi string) (Formula, bool) {

	firstRune, firstRuneSize := utf8.DecodeRune([]byte(phi))

	switch {
	case len(phi) == 0: // error
		return nil, false

	case firstRune == '□', firstRune == '◇', firstRune == '○', firstRune == '¬':
		if len(phi)+1 <= firstRuneSize || phi[firstRuneSize] != '(' {
			return nil, false
		}

		phi, ok := formulaFromStringRecursively(phi[firstRuneSize+1 : len(phi)-1])

		switch firstRune {
		case '□':
			return Always(phi), ok
		case '◇':
			return Eventually(phi), ok
		case '○':
			return Next(phi), ok
		case '¬':
			return Not(phi), ok
		}

	case firstRune == '(': // and / or / until
		bracketCount := 0
		for i := 1; i < len(phi); {
			b, runeSize := utf8.DecodeRune([]byte(phi[i:]))
			if b == '(' {
				bracketCount += 1
			} else if b == ')' {
				bracketCount -= 1
			} else if bracketCount == 0 && (b == '∧' || b == '∨') {
				subPhi, okPhi := formulaFromStringRecursively(phi[1:i])
				subPsi, okPsi := formulaFromStringRecursively(phi[i+runeSize : len(phi)-1])
				if b == '∧' {
					return And(subPhi, subPsi), okPhi && okPsi
				} else { // '∨'
					return Or(subPhi, subPsi), okPhi && okPsi
				}
			} else if bracketCount == 0 && b == 'U' { // until
				subPhi, okPhi := formulaFromStringRecursively(phi[2 : i-1])
				subPsi, okPsi := formulaFromStringRecursively(phi[i+2 : len(phi)-2])
				return Until(subPhi, subPsi), okPhi && okPsi
			} else if bracketCount == 0 && phi[i] == ')' {
				return nil, false
			}

			i += runeSize
		}

		return nil, false // too many opening brackets, or no ∧,∨ was found
	case firstRune != '¬' && firstRune != '(': // false, true, ap
		for i := 0; i < len(phi); {
			b, runeSize := utf8.DecodeRune([]byte(phi[i:]))
			if b == '¬' || b == '(' || b == ')' || b == '∧' || b == '∨' {
				return nil, false
			}
			i += runeSize
		}

		if phi == "false" {
			return False(), true
		} else if phi == "true" {
			return True(), true
		} else {
			return Ap(phi), true
		}
	}

	return nil, false // error
}
