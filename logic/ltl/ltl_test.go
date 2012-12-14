package ltl_test

import (
	"fmt"
	"github.com/hydroo/gomochex/logic/ltl"
	"testing"
)

func TestCreateFormula(t *testing.T) {

	if phi := ltl.Not(ltl.And(ltl.AP("π"), ltl.Not(ltl.AP("prΘp")))); fmt.Sprint(phi) != "¬((π∧¬(prΘp)))" {
		t.Error()
	}
}

func TestFormulaFromString(t *testing.T) {

	// correct formula
	if phi, ok := ltl.FormulaFromString("¬((π∧¬ (prΘp))    )"); ok != true || fmt.Sprint(phi) != "¬((π∧¬(prΘp)))" {
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
	if _, ok := ltl.FormulaFromString("((b))∧a)"); ok != false {
		t.Error()
	}

	// wrong brackets 3
	if _, ok := ltl.FormulaFromString("((b∧a)"); ok != false {
		t.Error()
	}
}

