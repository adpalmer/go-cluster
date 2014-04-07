/* NOTES: Needs to be optimized. Shouldn't be constantly moving around arrays.
 *        should instead be swapping pointers to those lists. That will be
 *        in next version.
 */

package gocluster

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type Cluster struct {
	Distance func(p1, p2 []float64) (float64, error)
}

// Function to randomly initialize clusters from input dataset.
// After initializing k clusters, runs Lloyd's Algorithm to find clusters
func (c Cluster) Km(entities [][]float64, k, maxIters int) ([][]float64, [][][]float64, error) {
	// create new random generator to avoid using the shared global object's Mutex lock
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numEntities := len(entities)

	// verify that all points have same dimensions
	if c.sameDimensionality(entities) == false {
		return nil, nil, errors.New("Dimensions must all be of same size")
	}

	// initialize entities
	centers := make([][]float64, k)
	for i := 0; i < k; i++ {
		centers[i] = entities[r.Intn(numEntities)]
	}
	return c.lloydsAlgo(entities, centers, maxIters)
}

func (c Cluster) Kmpp(entities [][]float64, k, maxIters int) ([][]float64, [][][]float64, error) {
	// create new random generator to avoid using the shared global object's Mutex lock
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numEntities := len(entities)

	// verify that all points have same dimensions
	if c.sameDimensionality(entities) == false {
		return nil, nil, errors.New("Dimensions must all be of same size")
	}

	centers := make([][]float64, k)
	d2 := make([]float64, numEntities) // distance squared
	centers[0] = entities[r.Intn(numEntities)]

	// iterate through rest of k to initialize other centers
	for i := 1; i < k; i++ {
		sum := 0.0
		for j := range entities {
			_, minDist := c.nearest(entities[j], centers[:i])
			d2[j] = minDist * minDist
			sum += d2[j]
		}

		// use find random number less than sum then iterate through d2 until index is found
		target := r.Float64() * sum
		j := 0
		for sum = d2[0]; sum < target; sum += d2[j] {
			j++
		}
		centers[i] = entities[j]
	}
	return c.lloydsAlgo(entities, centers, maxIters)
}

// Lloye's Algorithm for calculating k-means clustering
// called from Km
func (c Cluster) lloydsAlgo(entities, centers [][]float64, maxIters int) ([][]float64, [][][]float64, error) {
	numEntities := len(entities)

	// setup cluster groups
	clusters := make([][][]float64, len(centers))
	prevVariance := math.MaxFloat64

	// setup clusters
	// initial capacity of 10 may be too small. It will require doubling which
	// isn't ideal but will help lower overall memory if cluster sizes are small
	for i := range centers {
		clusters[i] = make([][]float64, 0, 10)
	}

	for i := 0; i < maxIters; i++ {
		// reset slices to zero
		for i := range centers {
			clusters[i] = clusters[i][:0]
		}

		curVariance := 0. // current variance
		for j := 0; j < numEntities; j++ {
			// get best cluster for each point
			minPt, minDist := c.nearest(entities[j], centers)
			clusters[minPt] = append(clusters[minPt], entities[j])
			curVariance += minDist
		}

		if prevVariance == curVariance {
			return centers, clusters, nil
		}

		// compute new centers
		for j := range clusters {
			centers[j] = c.updateCenter(clusters[j])
		}
		prevVariance = curVariance

	}
	return centers, clusters, nil
}

// update the cluster center given the points in the cluster
func (c Cluster) updateCenter(pts [][]float64) []float64 {
	if len(pts) == 0 {
		return make([]float64, 0)
	}

	dimensions := len(pts[0])
	numPts := len(pts)
	center := make([]float64, dimensions)

	for i := range pts {
		for j, v := range pts[i] {
			center[j] += v
		}
	}

	for i := range center {
		center[i] /= float64(numPts)
	}

	return center
}

// fine nearest center to given point
func (c Cluster) nearest(val []float64, centers [][]float64) (int, float64) {
	minDist, minPt := math.MaxFloat64, 0
	for l := range centers {
		// error caused when a cluster currently has no nodes (is a slice of size 0)
		dist, err := c.Distance(val, centers[l])
		if err == nil {
			if dist < minDist {
				minDist, minPt = dist, l
			}
		}
	}
	return minPt, minDist
}

// verify that pts are all of same size
func (c Cluster) sameDimensionality(pts [][]float64) bool {
	sz := len(pts[0])
	for i := range pts {
		if sz != len(pts[i]) {
			return false
		}
	}
	return true
}
