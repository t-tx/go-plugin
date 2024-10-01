package cplugin

import (
	"fmt"
	"sync/atomic"
	"time"
)

type config struct {
	discoverDirectory    string
	discoverInterval     time.Duration
	builtOutputDirectory string
	preloads             []*preload
}

func (c *config) SetDiscoverDirectory(discoverDirectory string) *config {
	c.discoverDirectory = discoverDirectory
	return c
}
func (c *config) SetDiscoverInterval(discoverInterval time.Duration) *config {
	c.discoverInterval = discoverInterval
	return c
}
func (c *config) SetBuiltOutputDirectory(builtOutputDirectory string) *config {
	c.builtOutputDirectory = builtOutputDirectory
	return c
}
func (c *config) Preload(preload ...*preload) *config {
	c.preloads = append(c.preloads, preload...)
	return c
}
func NewConfig() *config {
	return &config{}
}
func NewDefaultConfig() *config {
	return &config{
		discoverDirectory:    "resources/",
		discoverInterval:     2 * time.Second,
		builtOutputDirectory: "output/",
	}
}
func (c *config) getGoPath(name string) string {
	return name + ".go"
}

var counter atomic.Int32

func (c *config) genSoPath(name string) string {
	epochTime := time.Now().UnixNano()

	return fmt.Sprintf("%s_%d_%d.so", c.builtOutputDirectory+name, epochTime, counter.Add(1))
}

type preload struct {
	filePath string
	args     []string
	lazy     bool
}

func WithPreLoad(filePath string, lazy bool, args ...string) *preload {
	return &preload{
		filePath: filePath,
		args:     args,
		lazy:     lazy,
	}
}
