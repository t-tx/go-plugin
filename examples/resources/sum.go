package main

import (
	"github.com/t-tx/go-plugin/examples/types"
)

var Service types.Cal[int] = &speaker{}

type speaker struct {
}

func (s *speaker) Calculate(a ...int) int {
	var result int
	for _, v := range a {
		result = result + v
	}
	return result
}
