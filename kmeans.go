package kmeans

import (
	"math"
	"math/rand"
	"time"
)

type Cluster struct {
	Entity []interface{}
}

// Function to randomly initialize clusters from input dataset.
// After initializing k clusters, runs Lloyd's Algorithm to find clusters
func Km(calcDist func(p1, p2 interface{}) float64, updateClust func(Cluster) interface{}) func([]interface{}, int, int) ([]interface{}, []Cluster, error) {
	return func(entities []interface{}, k, maxIters int) ([]interface{}, []Cluster, error) {
		numEntities := len(entities)

		// initialize entities
		r := rand.New(rand.NewSource(time.Now().Unix()))
		centers := make([]interface{}, k)
		for i := 0; i < k; i++ {
			centers[i] = entities[r.Intn(numEntities)]
		}
		return lloydsAlgo(entities, centers, maxIters, calcDist, updateClust)
	}
}

// Lloye's Algorithm for calculating k-means clustering
// called from Km
func lloydsAlgo(entities, centers []interface{}, maxIters int, calcDist func(p1, p2 interface{}) float64, updateClust func(Cluster) interface{}) ([]interface{}, []Cluster, error) {
	numEntities := len(entities)

	// setup cluster groups
	clusters := make([]Cluster, len(centers))
	prevVariance := math.MaxFloat64

	for i := 0; i < maxIters; i++ {
		// setup clusters
		for i := range centers {
			clusters[i].Entity = make([]interface{}, 0, 10)
		}
		curVariance := 0. // current variance
		for j := 0; j < numEntities; j++ {
			// get best cluster for each point
			minDist, minPt := math.MaxFloat64, 0
			for l, val := range centers {
				dist := calcDist(entities[j], val)

				if dist < minDist {
					minDist, minPt = dist, l
				}
			}
			clusters[minPt].Entity = append(clusters[minPt].Entity, entities[j])
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
