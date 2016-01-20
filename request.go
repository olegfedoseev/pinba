// Package pinba provides convenient wrapper to work with
// pinba (https://github.com/tony2001/pinba_extension) protobuf messages
package pinba

import (
	proto "github.com/olegfedoseev/pinba/request"
	"strconv"
)

// Request wraps "raw" protobuf request and add convinient structs for Tags and Timers
type Request struct {
	proto.Request // protobuf Request

	Timers Timers
	Tags   Tags
}

// NewRequest unmarshal protobuf packet from given data,
// and convert raw tags and timers data to Tags and Timers structs
func NewRequest(data []byte) (*Request, error) {
	request := &Request{}
	if err := request.Unmarshal(data); err != nil {
		return nil, err
	}

	request.Tags = make(Tags, 4+len(request.TagValue))
	request.Tags[0] = Tag{Key: "host", Value: request.Hostname}
	request.Tags[1] = Tag{Key: "server", Value: request.ServerName}
	request.Tags[2] = Tag{Key: "script", Value: request.ScriptName}
	request.Tags[3] = Tag{Key: "status", Value: strconv.Itoa(int(request.Status))}

	for idx, val := range request.TagValue {
		request.Tags[4+idx] = Tag{
			Key:   request.Dictionary[request.TagName[idx]],
			Value: request.Dictionary[val],
		}
	}

	request.Timers = make([]Timer, len(request.TimerValue))
	offset := 0
	dictLen := uint32(len(request.Dictionary))

	for idx, val := range request.TimerValue {
		tags := make(Tags, int(request.TimerTagCount[idx])+4)
		tags[0] = Tag{Key: "host", Value: request.Hostname}
		tags[1] = Tag{Key: "server", Value: request.ServerName}
		tags[2] = Tag{Key: "script", Value: request.ScriptName}
		tags[3] = Tag{Key: "status", Value: strconv.Itoa(int(request.Status))}

		tagIdx := 4
		for valIdx, keyIdx := range request.TimerTagName[offset : offset+int(request.TimerTagCount[idx])] {
			if keyIdx >= dictLen || request.TimerTagValue[offset+valIdx] >= dictLen {
				continue
			}
			// Skip "global" tags
			if request.Dictionary[keyIdx] == "host" ||
				request.Dictionary[keyIdx] == "server" ||
				request.Dictionary[keyIdx] == "script" ||
				request.Dictionary[keyIdx] == "status" {
				continue
			}

			tags[tagIdx] = Tag{
				Key:   request.Dictionary[keyIdx],
				Value: request.Dictionary[request.TimerTagValue[offset+valIdx]],
			}
			tagIdx++
		}

		timer := Timer{
			Tags:     tags,
			HitCount: int32(request.TimerHitCount[idx]),
			Value:    val,
		}

		if len(request.TimerRuUtime) == len(request.TimerValue) {
			timer.RuUtime = request.TimerRuUtime[idx]
			timer.RuStime = request.TimerRuStime[idx]
		}
		request.Timers[idx] = timer
		offset += int(request.TimerTagCount[idx])
	}

	return request, nil
}
