package lib

import "github.com/samber/lo"

type Histogram interface {
	Process(data []int, height int) func(int) int
}

func ApplyHistogram(data []int, height int, histogram Histogram) []int {
	cb := histogram.Process(data, height)

	return lo.Map(data, func(v, _ int) int {
		return cb(v)
	})
}
