package kmeans

import (
	"math"
	"math/rand"
	"time"
)

type Cluster struct {
	Entity []interface{}
}

func LloydsAlgo(calcDist func(p1, p2 interface{}) float64, updateClust func(Cluster) interface{}) func([]interface{}, int, int) ([]interface{}, []Cluster, error) {
	return func(entities []interface{}, k, maxIters int) ([]interface{}, []Cluster, error) {
		numEntities := len(entities)

		// initialize entities
		r := rand.New(rand.NewSource(time.Now().Unix()))
		centers := make([]interface{}, k)
		for i := 0; i < k; i++ {
			centers[i] = entities[r.Intn(numEntities)]
		}

		// setup cluster groups
		clusters := make([]Cluster, k)
		prevVariance := math.MaxFloat64

		for i := 0; i < maxIters; i++ {
			// setup clusters
			for i := 0; i < k; i++ {
				clusters[i].Entity = make([]interface{}, 0, 10)
			}
			curVariance := 0. // current variance
			for j := 0; j < numEntities; j++ {
				// get best cluster for each point
				minDist, minPt := math.MaxFloat64, 0
				for l := 0; l < k; l++ {
					dist := calcDist(entities[j], centers[l])

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
			for j := 0; j < k; j++ {
				centers[j] = updateClust(clusters[j])
			}
			prevVariance = curVariance

		}
		return centers, clusters, nil
	}
}
