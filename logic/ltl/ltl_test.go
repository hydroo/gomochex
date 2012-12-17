package ltl_test

import (
	"fmt"
	"github.com/hydroo/gomochex/logic/ltl"
	"testing"
)

func TestCreateFormula(t *testing.T) {

	if phi := ltl.Or(ltl.Not(ltl.And(ltl.AP("π"), ltl.Not(ltl.AP("prΘp")))), ltl.AP("a")); fmt.Sprint(phi) != "(¬((π∧¬(prΘp)))∨a)" {
		t.Error()
	}
}

func TestFormulaFromString(t *testing.T) {

	// correct formula
	if phi, ok := ltl.FormulaFromString("(¬((π∧¬ (prΘp))    )∨a)"); ok != true || fmt.Sprint(phi) != "(¬((π∧¬(prΘp)))∨a)" {
		t.Error()
	}

	// wrong NOT
	if _, ok := ltl.FormulaFromString("¬x(a)"); ok != false {
		t.Error()
	}

	// correct AP
	if phi, ok := ltl.FormulaFromString("xx"); ok != true || fmt.Sprint(phi) != "xx" {
		t.Error()
	}

	// wrong AP
	if _, ok := ltl.FormulaFromString("x¬x"); ok != false {
		t.Error()
	}

	// wrong brackets 1
	if _, ok := ltl.FormulaFromString("(a∧(b)))"); ok != false {
		t.Error()
	}

	// wrong brackets 2
	if _, ok := ltl.FormulaFromString("((b))∨a)"); ok != false {
		t.Error()
	}

	// wrong brackets 3
	if _, ok := ltl.FormulaFromString("((b∧a)"); ok != false {
		t.Error()
	}

	// correct false
	if phi, ok := ltl.FormulaFromString("(¬(¬(false))∨false)"); ok != true || fmt.Sprint(phi) != "(¬(¬(false))∨false)" {
		t.Error()
	}

	// wrong false
	if phi, ok := ltl.FormulaFromString("(¬((falsefalse∧¬ (false))    )∨false)"); ok != true || fmt.Sprint(phi) != "(¬((falsefalse∧¬(false)))∨false)" {
		t.Error()
	}

	// correct true
	if phi, ok := ltl.FormulaFromString("(¬(¬(true))∨true)"); ok != true || fmt.Sprint(phi) != "(¬(¬(true))∨true)" {
		t.Error()
	}

	// wrong true
	if phi, ok := ltl.FormulaFromString("(¬((truetrue∧¬ (true))    )∨true)"); ok != true || fmt.Sprint(phi) != "(¬((truetrue∧¬(true)))∨true)" {
		t.Error()
	}

	// correct next
	if phi, ok := ltl.FormulaFromString("○(false)"); ok != true || fmt.Sprint(phi) != "○(false)" {
		fmt.Println(phi)
		t.Error()
	}

	// wrong next
	if _, ok := ltl.FormulaFromString("○false"); ok != false {
		t.Error()
	}

	// correct eventually
	if phi, ok := ltl.FormulaFromString("◇(false)"); ok != true || fmt.Sprint(phi) != "◇(false)" {
		fmt.Println(phi)
		t.Error()
	}

	// wrong eventually
	if _, ok := ltl.FormulaFromString("◇false"); ok != false {
		t.Error()
	}

	// correct until
	if phi, ok := ltl.FormulaFromString("((a)U(((b)U(false))))"); ok != true || fmt.Sprint(phi) != "((a)U(((b)U(false))))" {
		fmt.Println(phi)
		t.Error()
	}

}
