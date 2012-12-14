package xor

import (
	"fmt"
	"math/rand"
	"testing"
)

const K = 16

func TestXOR(t *testing.T) {
	m := &Metric{}
	for i := 0; i < K; i++ {
		m.Add(ID(i))
	}
	for piv := 0; piv < K; piv++ {
		nearest := m.Nearest(ID(piv), K/2)
		fmt.Println(ID(piv).ShortString(4))
		for _, n := range nearest {
			fmt.Println(" ", n.ID().ShortString(4))
		}
	}
}

const stressN = 1000000

func TestStress(t *testing.T) {
	m := &Metric{}
	var h []ID
	for i := 0; i < stressN; i++ {
		id := ID(rand.Int63())
		h = append(h, id)
		if _, err := m.Add(id); err != nil {
			t.Errorf("add (%s)", err)
		}
	}
	perm := rand.Perm(len(h))
	for _, j := range perm {
		m.Remove(h[j])
	}
}
