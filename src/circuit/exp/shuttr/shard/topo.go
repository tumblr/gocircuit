package shard

import (
	"circuit/exp/shuttr/x"
	"circuit/kit/xor"
)

type Shard struct {
	Pivot xor.Key
	Addr  x.Addr
	HTTP  int
}

func (sh *Shard) Key() xor.Key {
	return sh.Pivot
}

type Topo struct {
	metric xor.Metric
}

func New() *Topo {
	return &Topo{}
}

func NewPopulate(shards []*Shard) *Topo {
	t := &Topo{}
	t.Populate(shards)
	return t
}

func (t *Topo) Populate(shards []*Shard) {
	t.metric.Clear()
	for _, sh := range shards {
		t.Add(sh)
	}
}

func (t *Topo) Add(shard *Shard) {
	t.metric.Add(shard)
}

func (t *Topo) Find(key xor.Key) *Shard {
	nearest := t.metric.Nearest(key, 1)
	if len(nearest) == 0 {
		return nil
	}
	return nearest[0].(*Shard)
}

func (t *Topo) ChooseKey() xor.Key {
	return t.metric.ChooseMinK(5)
}
