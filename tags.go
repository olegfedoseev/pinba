package pinba

import (
	"bytes"
	"errors"
	"sort"
	"strings"
)

// MaxTags is how many tags we can have. OpenTSDB limits it to 8, so let it be 16 :)
const MaxTags = 16

// Tag is pair of key and value
type Tag struct {
	Key   string
	Value string
}

// Tags is a list of tags
type Tags []Tag

func (tags Tags) Len() int           { return len(tags) }
func (tags Tags) Swap(i, j int)      { tags[i], tags[j] = tags[j], tags[i] }
func (tags Tags) Less(i, j int) bool { return tags[i].Key < tags[j].Key }

// Get will return value of tag by given key or error if no such tag exists
func (tags Tags) Get(key string) (string, error) {
	for _, tag := range tags {
		if tag.Key == key {
			return tag.Value, nil
		}
	}
	return "", errors.New("no such tag")
}

// GetMap will tags as map
func (tags Tags) GetMap() map[string]string {
	result := map[string]string{}
	for _, tag := range tags {
		result[tag.Key] = tag.Value
	}
	return result
}

// Filter will filter (surprise :) tags by given keys, and return new Tags slice
// Code a bit strange, but that way we don't need to use append and it's almost
// twice as fast (480ns/op vs. 800ns/op)
func (tags Tags) Filter(filter []string) Tags {
	// First we count how many tags will be in result
	cnt := 0
	indexes := make([]int, MaxTags)
	for idx, tag := range tags {
		// Always skip empty tags
		if tag.Value == "" {
			continue
		}

		for _, f := range filter {
			if f == tag.Key {
				indexes[cnt] = idx
				cnt++
				break
			}
		}
	}

	// Then allocate only what we need, and loop once again to fill result
	result := make(Tags, cnt)
	for n, idx := range indexes[:cnt] {
		result[n] = tags[idx]
	}
	return result
}

// String will return tags as string such as "key=value key2=value2"
// Tags will be sorted by Key's
func (tags Tags) String() string {
	sort.Sort(tags)
	var buf bytes.Buffer
	for i, tag := range tags {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(tag.Key)
		buf.WriteString("=")
		buf.WriteString(tag.Value)
	}
	return buf.String()
}

// Stringf will return tags as formated string, with given format
// {xxx} will be replaced with value of tag xxx
func (tags Tags) Stringf(format string) string {
	result := format
	for _, tag := range tags {
		result = strings.Replace(result, "{"+tag.Key+"}", tag.Value, -1)
	}
	return result
}
