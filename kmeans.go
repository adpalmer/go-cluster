//package kmeans
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*   return func(entities []interface{}, k, maxIters int) ([]interface{}, []int, error) {
      numEntities := len(entities)
      fmt.Println(numEntities)
      return nil, nil, nil
   }
}*/

type Cluster struct {
	entity []interface{}
}

func kmeans(calcDist func(p1, p2 interface{}) float64, updateClust func(Cluster) interface{}) func([]interface{}, int, int) ([]interface{}, []Cluster, error) {
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
				clusters[i].entity = make([]interface{}, 0, 10)
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
				clusters[minPt].entity = append(clusters[minPt].entity, entities[j])
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
		fmt.Println(centers)
		return centers, clusters, nil
	}
}

type Point struct {
	x float64
	y float64
}

func euclideanDist(P1, P2 interface{}) float64 {
	p1, p2 := P1.(Point), P2.(Point)
	return math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
}

func updateCluster(c Cluster) interface{} {
	tot := Point{0., 0.}
	for i := 0; i < len(c.entity); i++ {
		tmp := c.entity[i].(Point)
		tot.x += tmp.x
		tot.y += tmp.y
	}
	tot.x /= float64(len(c.entity))
	tot.y /= float64(len(c.entity))
	return tot
}

func T(x []interface{}) {
	fmt.Println(x)
}

func main() {
	pts := []Point{Point{1, 1}, Point{2., .999}, Point{1, .3}, Point{10, 15}, Point{20, 10}, Point{15, 15}}

	// make it work with function
	interfacePoints := make([]interface{}, len(pts))

	for i, d := range pts {
		interfacePoints[i] = d
	}
	km := kmeans(euclideanDist, updateCluster)
	fmt.Println(km(interfacePoints, 2, 10))
	//kmeans(euclideanDist)
	//fmt.Println(euclideanDist(Point{1,1}, Point{1.5,2}))
}
