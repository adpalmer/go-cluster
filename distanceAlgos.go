package gocluster

import (
	"errors"
	"math"
)

func EuclideanDist(p1, p2 []float64) (float64, error) {
	if len(p1) != len(p2) {
		return math.MaxFloat64, errors.New("Points don't have the same number of dimensions")
	}

	sum := 0.
	var tmp float64
	for i := range p1 {
		tmp = p1[i] - p2[i]
		sum += tmp * tmp // should be quicker than using pow
	}
	return math.Sqrt(sum), nil
}
