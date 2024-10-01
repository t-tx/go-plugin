package main

import (
	"fmt"
	"time"

	"github.com/t-tx/go-plugin/examples/types"
	cplugin "github.com/t-tx/go-plugin/v1"
)

func main() {
	registry := cplugin.NewRegistry[*types.Cal[int]]()

	discovery := cplugin.NewDiscoveryWithRegistry(registry,
		cplugin.NewDefaultConfig().
			Preload(
				cplugin.WithPreLoad("resources/sum.go", false),
				cplugin.WithPreLoad("resources/min.go", false),
				cplugin.WithPreLoad("resources/max.go", false),
			))

	defer discovery.Destroy()

	for {
		func() {
			cal, ok := registry.Get("resources/sum")
			if !ok {
				fmt.Println("can't find sum")
				time.Sleep(time.Second * 2)
				return
			}
			fmt.Println("sum: ", (*cal).Calculate(1, 2, 3, 5, 6, 8, 9, 10))
			time.Sleep(time.Second * 2)
		}()
		func() {
			cal, ok := registry.Get("resources/max")
			if !ok {
				fmt.Println("can't find max")
				time.Sleep(time.Second * 2)
				return
			}
			fmt.Println("max: ", (*cal).Calculate(1, 2, 3, 5, 6, 8, 9, 10))
			time.Sleep(time.Second * 2)
		}()
		func() {
			cal, ok := registry.Get("resources/min")
			if !ok {
				fmt.Println("can't find min")
				time.Sleep(time.Second * 2)
				return
			}
			fmt.Println("min: ", (*cal).Calculate(1, 2, 3, 5, 6, 8, 9, 10))
			time.Sleep(time.Second * 2)
		}()
	}
}

type Cal interface {
	Sum(a, b int) int
}
