package intutil

import (
	"math"
	"strconv"
)

func Map(ints []int, fn func(int) int) []int {
	ret := make([]int, 0, len(ints))
	for _, i := range ints {
		ret = append(ret, fn(i))
	}
	return ret
}

func Max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func Min(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func ToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
