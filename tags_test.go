package pinba

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagsToSortedString(t *testing.T) {
	tags := Tags{Tag{"aaa", "val1"}, Tag{"bbb", "val2"}, Tag{"ccc", "val3"}}
	assert.Equal(t, "aaa=val1 bbb=val2 ccc=val3", tags.String())

	tags = Tags{Tag{"bbb", "val2"}, Tag{"aaa", "val1"}, Tag{"ccc", "val3"}}
	assert.Equal(t, "aaa=val1 bbb=val2 ccc=val3", tags.String())
}

func TestTagsGetTag(t *testing.T) {
	tags := Tags{Tag{"aaa", "val1"}, Tag{"bbb", "val2"}, Tag{"ccc", "val3"}}
	val, err := tags.Get("aaa")
	assert.Nil(t, err)
	assert.Equal(t, "val1", val)

	_, err = tags.Get("xxx")
	assert.NotNil(t, err)
}

func TestTagsFilter(t *testing.T) {
	tags := Tags{Tag{"aaa", "val1"}, Tag{"bbb", "val2"}, Tag{"ccc", "val3"}}
	assert.Equal(t, "aaa=val1 bbb=val2", tags.Filter([]string{"aaa", "bbb"}).String())

	tags = Tags{Tag{"bbb", "val2"}, Tag{"aaa", "val1"}, Tag{"ccc", "val3"}}
	assert.Equal(t, "aaa=val1 bbb=val2", tags.Filter([]string{"aaa", "bbb"}).String())

	tags = Tags{Tag{"bbb", "val2"}, Tag{"aaa", ""}, Tag{"ccc", "val3"}}
	assert.Equal(t, "bbb=val2", tags.Filter([]string{"aaa", "bbb"}).String())
}

func BenchmarkFilter(b *testing.B) {
	b.ResetTimer()

	tags := Tags{
		Tag{"key1", "val1"},
		Tag{"key2", ""},
		Tag{"key3", "val3"},
		Tag{"key4", "val4"},
		Tag{"key5", "val5"},
		Tag{"key6", ""},
		Tag{"key7", ""},
		Tag{"key8", "val8"},
		Tag{"key9", "val9"},
		Tag{"key10", "val10"},
	}

	filter := []string{"key1", "key3", "key5", "key6", "key9", "key10"}

	for i := 0; i < b.N; i++ {
		tags.Filter(filter)
	}
}
