package block

import (
	"time"
	"tumblr/app/sumr"
)

// Sketch summarizes a set of points. 
// Sketches can be combined to compute certain statistics and approximations.
type Sketch struct {
	UpdateTime time.Time // Application-level timestamp of the key
	Key        sumr.Key
	Sum        float64
	//Max      float64
	//Min      float64
	//SumSq    float64
	//SumAbs   float64
	//Count    uint32
}
