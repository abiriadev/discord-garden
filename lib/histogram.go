package lib

import (
	"slices"

	"github.com/samber/lo"
)

type Histogram interface {
	Process(data []int, height int) func(int) int
}

func ApplyHistogram(data []int, height int, histogram Histogram) []int {
	cb := histogram.Process(data, height)

	return lo.Map(data, func(v, _ int) int {
		return cb(v)
	})
}

type BinaryMeanHistogram struct{}

func (_ BinaryMeanHistogram) Process(data []int, height int) func(int) int {
	tmp := lo.Filter(data, func(item int, _ int) bool {
		return item > 0
	})
	l := len(tmp)
	slices.Sort(tmp)

	return func(v int) int {
		if v <= 0 {
			return v
		}
		i, _ := slices.BinarySearch(tmp, v)
		return i * height / l
	}
}
