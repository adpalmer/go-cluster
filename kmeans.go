package gocluster

import (
	"math"
	"math/rand"
)

// Function to randomly initialize clusters from input dataset.
// After initializing k clusters, runs Lloyd's Algorithm to find clusters
func Km(calcDist func(p1, p2 interface{}) float64, updateClust func([]interface{}) interface{}) func([]interface{}, int, int) ([]interface{}, [][]interface{}, error) {
	return func(entities []interface{}, k, maxIters int) ([]interface{}, [][]interface{}, error) {
		numEntities := len(entities)

		// initialize entities
		centers := make([]interface{}, k)
		for i := 0; i < k; i++ {
			centers[i] = entities[rand.Intn(numEntities)]
		}
		return lloydsAlgo(entities, centers, maxIters, calcDist, updateClust)
	}
}

func Kmpp(calcDist func(p1, p2 interface{}) float64, updateClust func([]interface{}) interface{}) func([]interface{}, int, int) ([]interface{}, [][]interface{}, error) {
	return func(entities []interface{}, k, maxIters int) ([]interface{}, [][]interface{}, error) {
		numEntities := len(entities)
		centers := make([]interface{}, k)
		d2 := make([]float64, numEntities) // distance squared
		centers[0] = entities[rand.Intn(numEntities)]

		// iterate through rest of k to initialize other centers
		for i := 1; i < k; i++ {
			sum := 0.0
			for j, data := range entities {
				_, minDist := nearest(data, centers[:i], calcDist)
				d2[j] = minDist * minDist
				sum += d2[j]
			}

			// use find random number less than sum then iterate through d2 until index is found
			target := rand.Float64() * sum
			j := 0
			for sum = d2[0]; sum < target; sum += d2[j] {
				j++
			}
			centers[i] = entities[j]
		}
		return lloydsAlgo(entities, centers, maxIters, calcDist, updateClust)
	}
}

// Lloye's Algorithm for calculating k-means clustering
// called from Km
func lloydsAlgo(entities, centers []interface{}, maxIters int, calcDist func(p1, p2 interface{}) float64, updateClust func([]interface{}) interface{}) ([]interface{}, [][]interface{}, error) {
	numEntities := len(entities)

	// setup cluster groups
	clusters := make([][]interface{}, len(centers))
	prevVariance := math.MaxFloat64

	for i := 0; i < maxIters; i++ {
		// setup clusters
		for i := range centers {
			clusters[i] = make([]interface{}, 0, 10)
		}
		curVariance := 0. // current variance
		for j := 0; j < numEntities; j++ {
			// get best cluster for each point
			minPt, minDist := nearest(entities[j], centers, calcDist)
			clusters[minPt] = append(clusters[minPt], entities[j])
			curVariance += minDist
		}

		if prevVariance == curVariance {
			return centers, clusters, nil
		}

		// compute new centers
		for j := range clusters {
			centers[j] = updateClust(clusters[j])
		}
		prevVariance = curVariance

	}
	return centers, clusters, nil
}

// fine nearest center to given point
func nearest(val interface{}, centers []interface{}, calcDist func(p1, p2 interface{}) float64) (int, float64) {
	minDist, minPt := math.MaxFloat64, 0
	for l, center := range centers {
		dist := calcDist(val, center)
		if dist < minDist {
			minDist, minPt = dist, l
		}
	}
	return minPt, minDist
}
