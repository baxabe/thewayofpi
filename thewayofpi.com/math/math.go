package math

import (
	"math"
	"errors"
)

const MaxInt = math.MaxInt64

func Max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func FindMin(list []int) (int, int, error) {
	var err error
	var idx, min int
	if len(list) > 0 {
		min = MaxInt
		for i, v := range list {
			if v < min {
				min = v
				idx = i
			}
		}
	} else {
		err = errors.New("math:FindMin - Empty input")
	}
	return min, idx, err
}

