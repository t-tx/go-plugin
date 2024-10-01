package main

import (
	"fmt"
	"time"

	"github.com/t-tx/go-plugin/examples/resources/def"
	"github.com/t-tx/go-plugin/examples/types"
	cplugin "github.com/t-tx/go-plugin/v1"
)

func main() {
	registry := cplugin.NewRegistry[*types.Cal[int]]()

	discovery := cplugin.NewDiscoveryWithRegistry(registry,
		cplugin.NewConfig().
			SetBuiltOutputDirectory("output/").
			SetDiscoverDirectory("resources/").
			SetDiscoverInterval(time.Second*1).
			Preload(
				cplugin.WithPreLoad("sum.go", false),
				cplugin.WithPreLoad("min.go", false),
				cplugin.WithPreLoad("max.go", false),
			))

	defer discovery.Destroy()

	for {
		func() {
			cal, ok := registry.Get("sum.go")
			if !ok {
				fmt.Println("can't find sum")
				time.Sleep(time.Second * 2)
				return
			}
			fmt.Println("sum: ", (*cal).Calculate(1, 2, 3, 5, 6, 8, 9, 10))
			time.Sleep(time.Second * 2)
		}()
		func() {
			cal, ok := registry.Get("max.go")
			if !ok {
				fmt.Println("can't find max")
				time.Sleep(time.Second * 2)
				return
			}
			fmt.Println("max: ", (*cal).Calculate(1, 2, 3, 5, 6, 8, 9, 10))
			time.Sleep(time.Second * 2)
		}()
		func() {
			cal, ok := registry.Get("min.go")
			if !ok {
				fmt.Println("can't find min")
				time.Sleep(time.Second * 2)
				return
			}
			fmt.Println("min: ", (*cal).Calculate(1, 2, 3, 5, 6, 8, 9, 10))
			time.Sleep(time.Second * 2)
		}()
		//can change the value of CAN_BE_ACCESS_FROM_PLUGIN in plugin => YES
		def.CAN_BE_ACCESS_FROM_PLUGIN++
	}
}

type Cal[T any] interface {
	Calculate(a ...T) T
}
