package main

import (
	"io/ioutil"
	"strconv"
	"tumblr/strings"
)

type Filter []int64

func ParseFilter(name string) (Filter, error) {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	tokens := strings.SplitTrim(string(raw), " \n\r\t")
	result := make(Filter, len(tokens))
	for i, t := range tokens {
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, err
		}
		result[i] = id
	}
	return result, nil
}
