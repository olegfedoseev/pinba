package pinba

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimersToString(t *testing.T) {
	tags := Tags{Tag{Key: "bla", Value: "foo1"}, Tag{Key: "foo", Value: "foo3"}}
	timers := Timers{
		Timer{Tags: tags, HitCount: 1, Value: 1.234, RuUtime: 0.123, RuStime: 0.012},
		Timer{Tags: tags, HitCount: 2, Value: 2.345, RuUtime: 0.234, RuStime: 0.023},
		Timer{Tags: tags, HitCount: 3, Value: 3.456, RuUtime: 0.345, RuStime: 0.034},
	}

	assert.Equal(t, timers.String(),
		"Tags: bla=foo1 foo=foo3, HitCount: 1, Value: 1.2340 Utime: 0.1230, Stime: 0.0120; "+
			"Tags: bla=foo1 foo=foo3, HitCount: 2, Value: 2.3450 Utime: 0.2340, Stime: 0.0230; "+
			"Tags: bla=foo1 foo=foo3, HitCount: 3, Value: 3.4560 Utime: 0.3450, Stime: 0.0340")
}
