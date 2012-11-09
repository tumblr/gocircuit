package stat

import (
	"time"
)

type MomentCarousel struct {
	duration  int64
	slots     []Moment
	head      int64
}

func NewMomentCarousel(nslots int, duration time.Duration) *MomentCarousel {
	x := &MomentCarousel{}
	x.Init(nslots, duration)
	return x
}

func (x *MomentCarousel) Init(nslots int, duration time.Duration) {
	slots := make([]Moment, nslots)
	for i, _ := range slots {
		slots[i].Init()
	}
	x.duration = int64(duration)
	x.slots = slots
}

func (x *MomentCarousel) TimeSpan() time.Duration {
	return time.Duration(x.duration * int64(len(x.slots)))
}

// Moment returns a pointer to the current moment structure corresponding to the time t
func (x *MomentCarousel) Slot(t time.Time) *Moment {
	slot := t.UnixNano() / x.duration
	if !x.spin(slot) {
		return nil
	}
	return &x.slots[int(slot % int64(len(x.slots)))]
}

func (x *MomentCarousel) Slots() ([]*Moment, time.Time) {
	result := make([]*Moment, len(x.slots))
	j := int(x.head % int64(len(x.slots))) + len(x.slots)
	for i := 0; i < len(result); i++ {
		result[i] = &x.slots[j % len(x.slots)]
		j--
	}
	return result, time.Unix(0, x.head * x.duration)
}

// spin rotates the circular slot buffer forward to ensure that the requested
// time falls within an interval slot. If the time t is before the earliest
// time in the buffer, spin is a nop and returns false.
func (x *MomentCarousel) spin(slot int64) bool {
	if slot + int64(len(x.slots)) <= x.head {
		return false
	}
	if slot <= x.head {
		return true
	}
	clear := int(min64(int64(len(x.slots)), slot - x.head))
	j := int((x.head + 1) % int64(len(x.slots)))
	for i := 0; i < clear; i++ {
		x.slots[j % len(x.slots)].Init()
		j++
	}
	x.head = slot
	return true
}

func (x *MomentCarousel) TailWeight(tail int) float64 {
	slots, _ := x.Slots()
	var result float64
	for i := 0; i < tail; i++ {
		result += float64(slots[i].Weight())
	}
	return result
}

func (x *MomentCarousel) Weight() float64 {
	var result float64
	for i, _ := range x.slots {
		result += x.slots[i].Weight()
	}
	return result
}

func (x *MomentCarousel) Mass() float64 {
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
