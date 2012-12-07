package a

import (
	_ "strings"
)

type T1 (int)

func (t T1) P1() {
	println(t)
}

func (t T1) P1() {
	println(t)
}

func (t *T1) P2() {
	println(t)
}

type (
	T2 struct {
	}

	T3 interface {
	}

	T4 []byte

	T5 [3]byte
)
