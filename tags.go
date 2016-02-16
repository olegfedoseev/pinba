package pinba

import (
	"bytes"
	"errors"
	"sort"
	"strings"
)

type Tag struct {
	Key   string
	Value string
}

type Tags []Tag

func (t Tags) Len() int           { return len(t) }
func (t Tags) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tags) Less(i, j int) bool { return t[i].Key < t[j].Key }

// Get will return value of tag by given key or error if no such tag exists
func (tags Tags) Get(key string) (string, error) {
	for _, tag := range tags {
		if tag.Key == key {
			return tag.Value, nil
		}
	}
	return "", errors.New("no such tag")
}

// Filter will filter (surprise :) tags by given keys, and return new Tags slice
// Code a bit strange, but that way we don't need to use append and it's almost
// twice as fast (480ns/op vs. 800ns/op)
func (tags Tags) Filter(filter []string) Tags {
	// First we count how many tags will be in result
	cnt := 0
	for _, tag := range tags {
		// Always skip empty tags
		if tag.Value == "" {
			continue
		}

		for _, f := range filter {
			if f == tag.Key {
				cnt++
				break
			}
		}
	}

	// Then allocate only what we need, and loop once again to fill result
	result := make(Tags, cnt)
	cnt = 0
	for _, tag := range tags {
		// Always skip empty tags
		if tag.Value == "" {
			continue
		}

		for _, f := range filter {
			if f == tag.Key {
				result[cnt] = tag
				cnt++
				break
			}
		}
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
