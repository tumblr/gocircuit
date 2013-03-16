// Package waterfill implements an algorithm for greedy even job assignment
package waterfill

import (
	"bytes"
	"fmt"
	"sort"
)

// Bin represents a bin carying integral load
type Bin interface {
	Add()
	Less(Bin) bool
}

// Fill represents an integral load distribution over a set of bins
type Fill struct {
	bin   []Bin
	i     int
	water Bin	// Bin holding the high water mark load
}

func NewFill(bin []Bin) *Fill {
	if len(bin) == 0 {
		return nil
	}
	sort.Sort(sortBins(bin))
	return &Fill{
		bin:   bin,
		i:     0,
		water: bin[0],
	}
}

func (f *Fill) String() string {
	var w bytes.Buffer
	for _, fb := range f.bin {
		s := fb.(fmt.Stringer)
		w.WriteString(s.String())
		w.WriteRune('·')
	}
	return string(w.Bytes())
}

// Add assigns a unit of work to a bin and returns that bin
func (f *Fill) Add() Bin {
	// Part I
	if f.i == len(f.bin) {
		f.i = 1
		r := f.bin[0]
		r.Add()
		f.water = r
		return r
	}
	// Part II
	r := f.bin[f.i]
	if r.Less(f.water) {
		r.Add()
		f.i++
		return r
	}
	// Part III
	f.i = 1
	r = f.bin[0]
	r.Add()
	f.water = r
	return r
}

// sortBins sorts a slice of Bins according to their order
type sortBins []Bin

func (sb sortBins) Len() int {
	return len(sb)
}

func (sb sortBins) Less(i, j int) bool {
	return sb[i].Less(sb[j])
}

func (sb sortBins) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}
