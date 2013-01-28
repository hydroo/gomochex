package ltl

import (
	"fmt"
	"testing"
)

func TestFormulaFromString(t *testing.T) {

	// correct formula
	if phi, ok := FormulaFromString("(¬((π∧¬ (prΘp))    )∨a)"); ok != true || phi.IsEqual(Or(Not(And(Ap("π"), Not(Ap("prΘp")))), Ap("a"))) != true {
		fmt.Println(phi)
		t.Error()
	}

	// wrong NOT
	if _, ok := FormulaFromString("¬x(a)"); ok != false {
		t.Error()
	}

	// correct AP
	if phi, ok := FormulaFromString("xx"); ok != true || phi.IsEqual(Ap("xx")) != true {
		t.Error()
	}

	// wrong AP
	if _, ok := FormulaFromString("x¬x"); ok != false {
		t.Error()
	}

	// wrong brackets 1
	if _, ok := FormulaFromString("(a∧(b)))"); ok != false {
		t.Error()
	}

	// wrong brackets 2
	if _, ok := FormulaFromString("((b))∨a)"); ok != false {
		t.Error()
	}

	// wrong brackets 3
	if _, ok := FormulaFromString("((b∧a)"); ok != false {
		t.Error()
	}

	// correct false
	if phi, ok := FormulaFromString("(¬(¬(false))∨false)"); ok != true || phi.IsEqual(Or(Not(Not(False())), False())) != true {
		t.Error()
	}

	// wrong false
	if phi, ok := FormulaFromString("(¬((falsefalse∧¬ (false))    )∨false)"); ok != true || phi.IsEqual(Or(Not(And(Ap("falsefalse"), Not(False()))), False())) != true {
		t.Error()
	}

	// correct true
	if phi, ok := FormulaFromString("(¬(¬(true))∨true)"); ok != true || phi.IsEqual(Or(Not(Not(True())), True())) != true {
		t.Error()
	}

	// wrong true
	if phi, ok := FormulaFromString("(¬((truetrue∧¬ (true))    )∨true)"); ok != true || phi.IsEqual(Or(Not(And(Ap("truetrue"), Not(True()))), True())) != true {
		t.Error()
	}

	// correct next
	if phi, ok := FormulaFromString("○(false)"); ok != true || phi.IsEqual(Next(False())) != true {
		fmt.Println(phi)
		t.Error()
	}

	// wrong next
	if _, ok := FormulaFromString("○false"); ok != false {
		t.Error()
	}

	// correct eventually
	if phi, ok := FormulaFromString("◇(false)"); ok != true || phi.IsEqual(Eventually(False())) != true {
		fmt.Println(phi)
		t.Error()
	}

	// wrong eventually
	if _, ok := FormulaFromString("◇false"); ok != false {
		t.Error()
	}

	// correct until
	if phi, ok := FormulaFromString("((a)U(((b)U(false))))"); ok != true || phi.IsEqual(Until(Ap("a"), Until(Ap("b"), False()))) != true {
		t.Error()
	}
}
