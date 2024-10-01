package main

import (
	"github.com/t-tx/go-plugin/examples/resources/def"
	"github.com/t-tx/go-plugin/examples/types"
)

var Service types.Cal[int] = &speaker{}

type speaker struct {
}

func (s *speaker) Calculate(a ...int) int {
	var min int = def.CAN_BE_ACCESS_FROM_PLUGIN
	for _, v := range a {
		if v < min || min == 0 {
			min = v
		}
	}
	return min
}
