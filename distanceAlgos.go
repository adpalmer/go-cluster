package gocluster

import (
  "math"
  "errors"
)

func EuclideanDist(p1, p2 []float64) (float64, error) {
  if len(p1) != len(p2) {
    return 0., errors.New("Points don't have the same number of dimensions")
  }

  sum := 0.
  var tmp float64
  for i := range p1 {
    tmp = p1[i] - p2[i]
    sum +=  tmp * tmp // should be quicker than using pow
  }
  return math.Sqrt(sum), nil
}
