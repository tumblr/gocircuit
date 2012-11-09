package block

import (
	"container/list"
	"fmt"
	"math"
	"sync"
	"time"
	"circuit/app/sumr"
	"circuit/kit/fs"
)

type Block struct {
	// Forget events older than forgetAfter
	forgetAfter time.Duration

	lk    sync.Mutex
	tabl  map[sumr.Key]*Sketch  // Key to sketch
	list  *list.List
	stat  Stat

	disk  *Disk
}

type Stat struct {
	NSketch      int64
	NWrite       int64
	NRead        int64
	NSketchInMem int64
}

func (s *Stat) String() string {
	return fmt.Sprintf("nsketch=%d nwrite=%d nread=%d", s.NSketch, s.NWrite, s.NRead)
}

func NewBlock(disk fs.FS, forgetAfter time.Duration) (*Block, error) {
	b := &Block{
		forgetAfter: forgetAfter,
		tabl:        make(map[sumr.Key]*Sketch),
		list:        list.New(),
	}
	if err := b.mountDisk(disk); err != nil {
		return nil, err
	}
	if forgetAfter > 0 {
		go b.forgetLoop()
	}
	return b, nil
}

func (b *Block) mountDisk(disk fs.FS) error {
	d, err := Mount(disk)	
	if err != nil {
		return err
	}
	file := d.Master()
	for {
		blob, err := file.Read()
		if err == ErrEndOfLog {
			break
		}
		if err != nil {
			return err
		}
		a, err := decodeAdd(blob)
		if err != nil {
			return err
		}
		b.Add(a.UTime, a.Key, a.Value)
	}
	b.disk = d
	return nil
}

// forgetLoop periodically deletes keys that have not been accessed for a specified amount of time
func (b *Block) forgetLoop() {
	for {
		time.Sleep(b.forgetAfter / 20)
		b.forget()
	}
}

// forget iterates through the keys and removes those not accessed in forgetAfter time.
func (b *Block) forget() {
	// Rationale: Database operations only add to the list of allocated
	// keys, therefore we can safely hold a pointer to the next iterate,
	// without holding a lock.  Releasing the lock on every iteration
	// allows competing database operations to share time equally.
	now := time.Now()
	b.lk.Lock()
	e := b.list.Front()
	b.lk.Unlock()
	for e != nil {
		b.lk.Lock()
		sketch := e.Value.(*Sketch)
		if _, present := b.tabl[sketch.Key]; !present {
			panic("sketch missing from table")
		}
		if now.Sub(sketch.UpdateTime) > b.forgetAfter {
			// This sketch is eligible to be expired.
			next := e.Next()
			b.list.Remove(e)
			delete(b.tabl, sketch.Key)
			b.stat.NSketchInMem--
			if b.stat.NSketchInMem < 0 {
				panic("sketches in memory less than 0")
			}
			e = next
		} else {
			e = e.Next()
		}
		b.lk.Unlock()
	}
}

// Add adds value to the current db value of key. If key does not exist, it is
// created with value equal to value. Add returns the "current" value of key,
// after the operation. Adding with value 0 can therefore be used to query the
// value of the key, without changing it.
func (b *Block) Add(updateTime time.Time, key sumr.Key, value float64) float64 {
	if math.IsNaN(value) {
		return math.NaN()
	}

	b.lk.Lock()
	defer b.lk.Unlock()

	if b.disk != nil {
		if _, err := b.disk.Master().Write(encodeAdd(&add{updateTime, key, value})); err != nil {
			panic(fmt.Sprintf("sumr block write (%s)", err))
		}
	}

	sketch := b.fetch(key)
	sketch.UpdateTime = maxTime(sketch.UpdateTime, updateTime)
	sketch.Sum += value
	b.stat.NWrite++
	return sketch.Sum
}

func (b *Block) Sum(key sumr.Key) float64 {
	b.lk.Lock()
	defer b.lk.Unlock()

	sketch := b.fetch(key)
	b.stat.NRead++
	return sketch.Sum
}

func (b *Block) fetch(key sumr.Key) *Sketch {
	sketch, present := b.tabl[key]
	if !present {
		sketch = &Sketch{Key: key}
		b.tabl[key] = sketch
		b.list.PushBack(sketch)
		b.stat.NSketchInMem++
		b.stat.NSketch++
	}
	return sketch
}

// Stat returns an object containing statistics about the operation of this database block.
func (b *Block) Stat() *Stat {
	b.lk.Lock()
	defer b.lk.Unlock()
	r := b.stat
	return &r
}

func maxTime(p, q time.Time) time.Time {
	if p.Sub(q) >= 0 {
		return p
	}
	return q
}
