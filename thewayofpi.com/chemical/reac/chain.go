package reac

import (
	"errors"
	"fmt"
	"math"
	"thewayofpi.com/buffer/random"
)

func buildStepsChain(start, end, steps uint) ([]uint, error) {
	if !chainExists(start, end, steps) {
		return nil, errors.New("reac:buildStepsChain - no path")
	}
	if steps == 1 {
		return []uint{start, end}, nil
	}
	seq := make([]uint, steps+1)
	seq[0] = start
	seq[steps] = end
	if err := stepsChainBuilder(seq); err != nil {
		return nil, errors.New(fmt.Sprintf("reac:buildStepsChainRec - %v\n%v\n", err, seq))
	}
	return seq[:], nil
}

func stepsChainBuilder(path []uint) error {
	rnd := random.New()
	for first, last := 0, len(path)-1; last-first > 2; first, last = first+1, last-1 {
		if last-first == 3 {
			// Check 2**2
			if path[first] == 2 && path[last] == 2 {
				path[first+1] = uint(rnd.Intn(2) + 3)
				path[last-1] = uint(rnd.Intn(2) + 3)
				break
			}
			// Check 1**1
			if path[first] == 1 && path[last] == 1 {
				var opt = [...]int{0, 3}
				i := rnd.Intn(len(opt))
				path[first+opt[i]] = uint(rnd.Intn(2) + 3)
			}
		}
		if val, err := nextStep(path[first], path[last], uint(last-first-1)); err != nil {
			return errors.New(fmt.Sprintf("reac:stepsChainBuilder - %v", err))
		} else {
			path[first+1] = val
		}
		if val, err := nextStep(path[last], path[first+1], uint(last-first-2)); err != nil {
			return errors.New(fmt.Sprintf("reac:stepsChainBuilder - %v", err))
		} else {
			path[last-1] = val
		}
		// Check *22*
		if last-first == 3 && path[first+1] == 2 && path[last-1] == 2 {
			if path[first] != 1 {
				path[first+1] = uint(rnd.Intn(2) + 3)
			} else if path[last] != 1 {
				path[last-1] = uint(rnd.Intn(2) + 3)
			} else {
				return errors.New(fmt.Sprintf("reac:stepsChainBuilder - Check *22* failed [%v22%v]", path[0], path[3]))
			}
		}
	}
	if len(path)%2 == 1 {
		center := len(path) / 2
		// Check 2*1 & 1*2
		if path[center-1]+path[center+1] == 3 {
			path[center] = uint(rnd.Intn(2) + 3)
			if path[center-1] == 1 {
				path[center-1] = uint(rnd.Intn(2) + 3)
			} else {
				path[center+1] = uint(rnd.Intn(2) + 3)
			}
		} else {
			if val, err := nextStep(path[center-1], path[center+1], 1); err != nil {
				return errors.New(fmt.Sprintf("reac:stepsChainBuilder - %v", err))
			} else {
				if val == 2 && (path[center-1] == 2 || path[center+1] == 2) {
					val = uint(rnd.Intn(2) + 3)
				}
				path[center] = val
			}
		}
	}
	return nil
}

func chainExists(pointA, pointB, steps uint) bool {
	if steps == 0 {
		return false
	}
	if steps == 1 && pointA == pointB && (pointA == 1 || pointA == 2) {
		return false
	}
	if steps == 2 && pointA+pointB == 3 {
		return false
	}
	if steps == 3 && pointA == 1 && pointB == 1 {
		return false
	}
	maxA, minA := calcMaxMin(pointA, steps)
	maxB, minB := calcMaxMin(pointB, steps)
	if pointB >= minA && pointB <= maxA && pointA >= minB && pointA <= maxB {
		return true
	}
	return false
}

func calcMaxMin(point, dist uint) (max, min uint) {
	factor := math.Exp2(float64(dist))
	max = uint(float64(point) * factor)
	min = uint(math.Max(1, math.Ceil(float64(point)/factor)))
	return
}

func nextStep(start, end, dist uint) (uint, error) {
	if dist == 0 {
		return 0, errors.New("reac:nextStep - no steps to do")
	}
	rnd := random.New()
	startMax := 2 * start
	startMin := uint(start/2 + start%2)
	endMax, endMin := calcMaxMin(end, dist)
	max := uint(math.Min(float64(startMax), float64(endMax)))
	min := uint(math.Max(float64(startMin), float64(endMin)))
	if max < min {
		return 0, errors.New(fmt.Sprintf("reac:nextStep - no path: %v -> %v (%v) [%v..%v]", start, end, dist, min, max))
	}
	val := min + uint(rnd.Intn(int(max-min)+1))
	if start == 1 && val == 1 {
		val++
	} else if start == 2 && val == 2 {
		var opt = [...]uint{3, 4}
		i := rnd.Intn(len(opt))
		val = opt[i]
	}
	return val, nil
}

