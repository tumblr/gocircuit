package a

type T1 int

func (t T1) P1() {
	println(t)
}

type (
	T2 struct {
	}

	T3 interface {
	}
)
