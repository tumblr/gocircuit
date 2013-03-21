package stat

import (
	"time"
)

type SlidingMoment struct {
	slotdur int64
	slots   []Moment
	head    int64
}

func NewSlidingMoment(resolution int, duration time.Duration) *SlidingMoment {
	x := &SlidingMoment{}
	x.Init(resolution, duration)
	return x
}

func (x *SlidingMoment) Init(resolution int, duration time.Duration) {
	slots := make([]Moment, resolution)
	for i, _ := range slots {
		slots[i].Init()
	}
	x.slotdur = int64(duration) / int64(resolution)
	x.slots = slots
}

func (x *SlidingMoment) TimeSpan() time.Duration {
	return time.Duration(x.slotdur * int64(len(x.slots)))
}

// Moment returns a pointer to the current moment structure corresponding to the time t
func (x *SlidingMoment) Slot(t time.Time) *Moment {
	slot := t.UnixNano() / x.slotdur
	if !x.spin(slot) {
		return nil
	}
	return &x.slots[int(slot%int64(len(x.slots)))]
}

func (x *SlidingMoment) Slots() ([]*Moment, time.Time) {
	result := make([]*Moment, len(x.slots))
	j := int(x.head%int64(len(x.slots))) + len(x.slots)
	for i := 0; i < len(result); i++ {
		result[i] = &x.slots[j%len(x.slots)]
		j--
	}
	return result, time.Unix(0, x.head*x.slotdur)
}

// spin rotates the circular slot buffer forward to ensure that the requested
// time falls within an interval slot. If the time t is before the earliest
// time in the buffer, spin is a nop and returns false.
func (x *SlidingMoment) spin(slot int64) bool {
	if slot+int64(len(x.slots)) <= x.head {
		return false
	}
	if slot <= x.head {
		return true
	}
	clear := int(min64(int64(len(x.slots)), slot-x.head))
	j := int((x.head + 1) % int64(len(x.slots)))
	for i := 0; i < clear; i++ {
		x.slots[j%len(x.slots)].Init()
		j++
	}
	x.head = slot
	return true
}

func (x *SlidingMoment) TailWeight(tail int) float64 {
	slots, _ := x.Slots()
	var result float64
	for i := 0; i < tail; i++ {
		result += float64(slots[i].Weight())
	}
	return result
}

func (x *SlidingMoment) Weight() float64 {
	var result float64
	for i, _ := range x.slots {
		result += x.slots[i].Weight()
	}
	return result
}

func (x *SlidingMoment) Mass() float64 {
	var result float64
	for i, _ := range x.slots {
		result += x.slots[i].Mass()
	}
	return result
}

func min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
