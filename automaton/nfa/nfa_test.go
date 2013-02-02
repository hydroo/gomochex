package nfa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hydroo/gomochex/basic/set"
	"testing"
)

//                           +-+
// --> 1 -- a --> 2 -- b --> |3|
//     ^                     +-+
//     |                      |
//     +--------- a ----------+
//
func TestSimpleNfa(t *testing.T) {
	a := Letter("a")
	b := Letter("b")
	q0 := State("0")
	q1 := State("1")
	q2 := State("2")

	A := NewNfa()
	A.Alphabet().Add(a, b)
	A.States().Add(q0, q1, q2)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(q2)

	trans := func(s State, l Letter) StateSet {
		if s.IsEqual(q0) && l.IsEqual(a) {
			return set.NewSet(q1)
		} else if s.IsEqual(q1) && l.IsEqual(b) {
			return set.NewSet(q2)
		} else if s.IsEqual(q2) && l.IsEqual(a) {
			return set.NewSet(q0)
		}
		return set.NewSet()
	}
	A.SetTransitionFunction(trans)

	s := []byte(`{"States":["0","1","2"],"Alphabet":["a","b"],"InitialStates":["0"],"Transitions":{"0":{"a":["1"]},"1":{"b":["2"]},"2":{"a":["0"]}},"FinalStates":["2"]}`)
	compareNfaToMarshaledNfa(A, s, true, t, "")
}

func TestConcat(t *testing.T) {
	A := Concat(OneLetter("a"), Union(OneLetter("π"), OneLetter("c"))).(*simpleNfa).removeUselessParts()

	//       -- a --> o -- π --> □
	//      /
	// --> o
	//      \
	//       -- a --> o -- c --> □
	s := []byte(`{"States":["0","1","2","3","4"],"Alphabet":["a","π","c"],"InitialStates":["0"],"Transitions":{"0":{"a":["1","3"]},"1":{"π":["2"]},"3":{"c":["4"]}},"FinalStates":["2","4"]}`)
	compareNfaToMarshaledNfa(A, s, true, t, "")
}

func TestConcatNfa2(t *testing.T) {
	a := Letter("a")
	A := Concat(OneLetter(a), OneLetter(a)).(*simpleNfa).removeUselessParts()

	s := []byte(`{"States":["0","1","2"],"Alphabet":["a"],"InitialStates":["0"],"Transitions":{"0":{"a":["1"]},"1":{"a":["2"]}},"FinalStates":["2"]}`)
	compareNfaToMarshaledNfa(A, s, true, t, "")
}

func TestKleeneStar(t *testing.T) {
	a := Letter("a")
	b := Letter("π")
	A := KleeneStar(Union(Concat(OneLetter(a), OneLetter(a)), Concat(OneLetter(b), OneLetter(b)))).(*simpleNfa).removeUselessParts()

	// minimal nfa:
	//
	//   +-- a ---
	//   |        \
	//   | - a -> o
	//   v/
	// ->□
	//   ^\
	//   | - π -> o
	//   |       /
	//   +-- π --
	s := []byte(`{"States":["0","1","2"],"Alphabet":["a","π"],"InitialStates":["0"],"Transitions":{"0":{"a":["1"],"π":["2"]},"1":{"a":["0"]},"2":{"π":["0"]}},"FinalStates":["0"]}`)
	compareNfaToMarshaledNfa(A, s, true, t, "")
}

func TestOneLetter(t *testing.T) {
	a := Letter("π")
	A := OneLetter(a)

	s := []byte(`{"States":["0","1"],"Alphabet":["π"],"InitialStates":["0"],"Transitions":{"0":{"π":["1"]}},"FinalStates":["1"]}`)
	compareNfaToMarshaledNfa(A, s, true, t, "")
}

func TestUnion(t *testing.T) {
	a := Letter("a")
	b := Letter("π")
	A := Union(OneLetter(a), OneLetter(b))

	s := []byte(`{"States":["0","1","2","3"],"Alphabet":["a","π"],"InitialStates":["0","2"],"Transitions":{"0":{"a":["1"]},"2":{"π":["3"]}},"FinalStates":["1","3"]}`)

	compareNfaToMarshaledNfa(A, s, true, t, "")
}

func TestJson(t *testing.T) {
	a := Letter("a")
	q0 := State("0")
	q1 := State("1")
	q2 := State("2")
	q3 := State("3")

	// -> □ -> o
	//
	//    o -> □
	A := NewNfa()
	A.Alphabet().Add(a)
	A.States().Add(q0, q1, q2, q3)
	A.InitialStates().Add(q0)
	A.FinalStates().Add(q0, q3)
	A.SetTransition(q0, a, set.NewSet(q1))
	A.SetTransition(q2, a, set.NewSet(q3))

	s1, err1 := json.Marshal(A)
	if err1 != nil {
		t.Error(err1)
	}

	B := NewNfa()
	err2 := json.Unmarshal(s1, &B)
	if err2 != nil {
		t.Error(err2)
	}

	s3, err3 := json.Marshal(B)
	if err1 != nil {
		t.Error(err3)
	}

	if bytes.Compare(s1, s3) != 0 {
		t.Error()
	}
}

func TestIsEqual(t *testing.T) {
	type test struct {
		s, u          []byte
		shouldBeEqual bool
	}

	tests := []test{
		test{
			[]byte(`{"States":["0","1"],"Alphabet":["a"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["0"]},"1":{"a":["1"]}},"FinalStates":["0","1"]}`),
			[]byte(`{"States":["1","0"],"Alphabet":["a"],"InitialStates":["1","0"],"Transitions":{"1":{"a":["1"]},"0":{"a":["0"]}},"FinalStates":["0","1"]}`),
			true},
		test{
			[]byte(`{"States":["0","1"],"Alphabet":["a"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["0"]},"1":{"a":["1"]}},"FinalStates":["0","1"]}`),
			[]byte(`{"States":["x","y"],"Alphabet":["a"],"InitialStates":["x","y"],"Transitions":{"x":{"a":["x"]},"y":{"a":["y"]}},"FinalStates":["x","y"]}`),
			true},
		test{
			[]byte(`{"States":["0","1"],"Alphabet":["a"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["0"]},"1":{"a":["1"]}},"FinalStates":["0","1"]}`),
			[]byte(`{"States":["0","1","2"],"Alphabet":["a"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["0"]},"1":{"a":["1"]}},"FinalStates":["0","1"]}`),
			false},
		test{
			[]byte(`{"States":["0","1"],"Alphabet":["a"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["0"]},"1":{"a":["1"]}},"FinalStates":["0","1"]}`),
			[]byte(`{"States":["0","2"],"Alphabet":["a"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["0"]},"1":{"a":["1"]}},"FinalStates":["0","1"]}`),
			false},
		test{
			[]byte(`{"States":["0"],"Alphabet":["a"],"InitialStates":["0"],"Transitions":{},"FinalStates":["0"]}`),
			[]byte(`{"States":["0"],"Alphabet":["a"],"InitialStates":["0"],"Transitions":{"0":{"a":["0"]}},"FinalStates":["0"]}`),
			false},
		test{
			[]byte(`{"States":["0"],"Alphabet":["a","b"],"InitialStates":["0"],"Transitions":{"0":{"a":["0"]}},"FinalStates":["0"]}`),
			[]byte(`{"States":["0"],"Alphabet":["a","b"],"InitialStates":["0"],"Transitions":{"0":{"b":["0"]}},"FinalStates":["0"]}`),
			false},
		test{
			[]byte(`{"States":["0","1","2","3"],"Alphabet":["a","π"],"InitialStates":["0","2"],"Transitions":{"0":{"a":["1"]},"2":{"π":["3"]}},"FinalStates":["1","3"]}`),
			[]byte(`{"States":["0","1","2","3"],"Alphabet":["a","π"],"InitialStates":["0","1"],"Transitions":{"0":{"a":["1"]},"2":{"π":["3"]}},"FinalStates":["1","3"]}`),
			false},
	}

	for k, x := range tests {
		compareMarshaledNfaToMarshaledNfa(x.s, x.u, x.shouldBeEqual, t, fmt.Sprint("case ", k))
	}
}

func TestInducedNfa(t *testing.T) {
	// -> o -> o -> □
	//    |    ^    |
	//    |    +----
	//    |
	//    + -> □
	s := []byte(`{"States":["0","1","2","3"],"Alphabet":["a"],"InitialStates":["0"],"Transitions":{"0":{"a":["1","3"]},"1":{"a":["2"]},"2":{"a":["1"]}},"FinalStates":["2","3"]}`)

	A := NewNfa().(*simpleNfa)
	json.Unmarshal(s, &A)

	results := make([]struct {
		should, is []byte
		err        error
	}, 4)

	results[0].is, results[0].err = json.Marshal(A.inducedNfa(State("0")))
	results[1].is, results[1].err = json.Marshal(A.inducedNfa(State("1")))
	results[2].is, results[2].err = json.Marshal(A.inducedNfa(State("3")))
	results[3].is, results[3].err = json.Marshal(A.inducedNfa(State("2")))

	results[0].should = s
	results[1].should = []byte(`{"States":["1","2"],"Alphabet":["a"],"InitialStates":["1"],"Transitions":{"1":{"a":["2"]},"2":{"a":["1"]}},"FinalStates":["2"]}`)
	results[2].should = []byte(`{"States":["3"],"Alphabet":["a"],"InitialStates":["3"],"Transitions":{},"FinalStates":["3"]}`)
	results[3].should = []byte(`{"States":["2","1"],"Alphabet":["a"],"InitialStates":["2"],"Transitions":{"1":{"a":["2"]},"2":{"a":["1"]}},"FinalStates":["2"]}`)

	for k, r := range results {
		if r.err != nil || bytes.Compare(r.is, r.should) != 0 {
			t.Error("case ", k, "\nshould:", string(r.should), "\nis:    ", string(r.is))
		}
	}
}

func TestRemoveUselessParts(t *testing.T) {
	//    o
	//    |
	//    v
	// -> □ -> o
	//
	//    o -> □
	//
	// -> o -> o
	//
	// -> o -> □
	s := []byte(`{
		"States":["0","1","2","3","4","5","6","7","8"],
		"Alphabet":["a"],
		"InitialStates":["0","4","7"],
		"Transitions":{"0":{"a":["1"]},"2":{"a":["3"]},"4":{"a":["5"]},"6":{"a":["0"]},"7":{"a":["8"]}},
		"FinalStates":["0","3","8"]
	}`)
	u := []byte(`{"States":["0","7","8"],"Alphabet":["a"],"InitialStates":["0","7"],"Transitions":{"7":{"a":["8"]}},"FinalStates":["0","8"]}`)

	A := NewNfa().(*simpleNfa)
	json.Unmarshal(s, &A)
	B := A.removeUselessParts()
	compareNfaToMarshaledNfa(B, u, true, t, "")
}

func compareMarshaledNfaToMarshaledNfa(s, u []byte, shouldBeEqual bool, t *testing.T, extraInfo string) {
	A := NewNfa()
	json.Unmarshal(s, &A)
	compareNfaToMarshaledNfa(A, u, shouldBeEqual, t, extraInfo)
}

func compareNfaToMarshaledNfa(A Nfa, s []byte, shouldBeEqual bool, t *testing.T, extraInfo string) {
	B := NewNfa()
	json.Unmarshal(s, &B)
	if A.IsEqual(B) != shouldBeEqual {
		u, _ := json.Marshal(A)
		if extraInfo != "" {
			extraInfo = fmt.Sprint("\n", extraInfo)
		}
		t.Error(extraInfo, "\nshould be equal: ", shouldBeEqual, " \nshould:", string(s), "\nis:    ", string(u))
	}
}
