package lib

type Histogram interface {
	Process(data []int, height int) func(int) int
}
