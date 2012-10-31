// Package stat implements facilities for storing distribution sketches and updating them
package stat

import (
	"math"
)

type Moment struct {
	sum	float64
	sumAbs	float64
	sumSq	float64
	min     float64
	max     float64
	weight  float64
	n       int64
}

func (x *Moment) Init() {
	x.sum, x.sumAbs, x.sumSq, x.min, x.max, x.weight, x.n = 0, 0, 0, math.NaN(), math.NaN(), 0, 0
}

func (x *Moment) Add(sample float64) {
	x.AddWeighted(sample, 1)
}

func (x *Moment) AddWeighted(sample float64, weight float64) {
	x.n++
	x.sum += sample*weight
	x.sumAbs += math.Abs(sample*weight)
	x.sumSq += sample*sample*weight
	x.weight += weight
	if math.IsNaN(x.min) || sample < x.min {
		x.min = sample
	}
	if math.IsNaN(x.max) || sample > x.max {
		x.max = sample
	}
}

func (x *Moment) IsEmpty() bool {
	return x.n == 0
}

func (x *Moment) Count() int64 {
	return x.n
}

func (x *Moment) Weight() float64 {
	return x.weight
}

func (x *Moment) Mass() float64 {
	return x.sumAbs
}

func (x *Moment) Average() float64 {
	return x.Moment(1)
}

func (x *Moment) Variance() float64 {
	m1 := x.Moment(1)
	return x.Moment(2) - m1*m1
}

func (x *Moment) StdDev() float64 {
	return math.Sqrt(x.Variance())
}

func (x *Moment) Min() float64 {
	return x.min
}

func (x *Moment) Max() float64 {
	return x.max
}

func (x *Moment) Moment(k float64) float64 {
	switch k {
	case 0:
		return 1
	case 1:
		return x.sumAbs / x.weight
	case 2:
		return x.sumSq / x.weight
	}
	if math.IsInf(k, 1) {
		return x.max
	}
	panic("not yet supported")
}
