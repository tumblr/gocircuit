package xor

import "math/rand"

// ChooseMinK tries to insert K random IDs ...
func (m *Metric) ChooseMinK(k int) ID {
	if m == nil {
		return ID(rand.Int63())
	}
	var min_id ID
	var min_d int = 1000
	for k > 0 {
		// Note: The last bit is not really randomized here
		id := ID(rand.Int63())
		d, err := m.Add(id)
		if err != nil {
			continue
		}
		m.Remove(id)
		if d < min_d {
			min_id = id
			min_d = d
		}
		k--
	}
	return min_id
}
