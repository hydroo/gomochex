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
