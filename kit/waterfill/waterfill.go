package waterfill

import (
	"bytes"
	"fmt"
	"sort"
)

// FillBin represents a bin carying integer load
type FillBin interface {
	Add()
	Less(FillBin) bool
}

func NewFill(bin []FillBin) *Fill {
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

type Fill struct {
	bin   []FillBin
	i     int
	water FillBin
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

func (f *Fill) Add() FillBin {
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

// sortBins sorts a slice of FillBins according to their order
type sortBins []FillBin

func (sb sortBins) Len() int {
	return len(sb)
}

func (sb sortBins) Less(i, j int) bool {
	return sb[i].Less(sb[j])
}

func (sb sortBins) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}
