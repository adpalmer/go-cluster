package main

import (
	"fmt"
	"github.com/adpalmer/go-cluster"
)

func main() {
	pts := [][]float64{[]float64{1, 1}, []float64{2., .999}, []float64{1, .3}, []float64{10, 15}, []float64{20, 10}, []float64{15, 15}}

	// Solve with standard k-means
	cluster := gocluster.Cluster{gocluster.EuclideanDist}
	clusterCenters, clusterMembers, _ := cluster.Km(pts, 2, 10)
	for i := 0; i < len(clusterCenters); i++ {
		fmt.Printf("Cluster Center %v:\n", clusterCenters[i])
		fmt.Println("\tmembers -> ", clusterMembers[i])
	}

	// Solve with k-means++
	cluster = gocluster.Cluster{gocluster.EuclideanDist}
	clusterCenters, clusterMembers, _ = cluster.Kmpp(pts, 2, 10)
	for i := 0; i < len(clusterCenters); i++ {
		fmt.Printf("Cluster Center %v:\n", clusterCenters[i])
		fmt.Println("\tmembers -> ", clusterMembers[i])
	}
}
