package set

type Element interface {
	IsEqual(e Element) bool
}

type Copier interface {
	Copy() Copier
}

type Newer interface {
	New() Newer
}

type Set interface {
	Copier
	Newer
	Add(e Element)
	At(index int) (Element, bool)
	Probe(e Element) bool
	Remove(e Element)
	Size() int
}

func Intersect(S, T Set) Set {
	ret := S.New().(Set)

	for i := 0; i < S.Size(); i += 1 {
		for j := 0; j < T.Size(); j += 1 {
			s, _ := S.At(i)
			t, _ := T.At(j)
			if s.IsEqual(t) {
				ret.Add(s)
				break
			}
		}
	}

	return ret
}

func Join(S, T Set) Set {
	ret := S.Copy().(Set)

	for i := 0; i < T.Size(); i += 1 {
		e, _ := T.At(i)
		ret.Add(e)
	}

	return ret
}


type SimpleSet []Element

func NewSimpleSet() *SimpleSet {
	return new(SimpleSet)
}

func (S *SimpleSet) Add(e Element) {
	if S.Probe(e.(Element)) == false {
		*S = append(*S, e.(Element))
	}
}

func (S SimpleSet) At(index int) (Element, bool) {
	if index < len(S) {
		return S[index], true
	} //else {
	return S[0], false
	//}
}

func (S SimpleSet) Copy() Copier {
	cp := make(SimpleSet, S.Size())
	copy(cp, S)
	return &cp
}

func (S SimpleSet) New() Newer {
	return new(SimpleSet)
}

func (S SimpleSet) Probe(e Element) bool {
	for _, v := range S {
		if e.IsEqual(v) {
			return true
		}
	}

	return false
}

func (S *SimpleSet) Remove(e Element) {

	if S.Probe(e) == false {
		return
	}

	cp := new(SimpleSet)

	for _, v := range *S {
		if e.IsEqual(v) == false {
			cp.Add(v)
		}
	}

	*S = *cp
}

func (S SimpleSet) Size() int {
	return len(S)
}

