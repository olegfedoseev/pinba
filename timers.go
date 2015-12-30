package pinba

import (
	"bytes"
	"fmt"
)

type Timer struct {
	Tags     Tags
	HitCount int32
	Value    float32
	RuUtime  float32
	RuStime  float32
}

type Timers []Timer

func (t Timers) String() string {
	b := &bytes.Buffer{}
	for i, timer := range t {
		if i > 0 {
			fmt.Fprint(b, "; ")
		}
		fmt.Fprintf(b, "Tags: %s, HitCount: %d, Value: %3.4f Utime: %3.4f, Stime: %3.4f",
			timer.Tags.String(), timer.HitCount, timer.Value, timer.RuUtime, timer.RuStime)
	}
	return b.String()
}
