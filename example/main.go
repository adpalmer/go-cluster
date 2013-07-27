package main

import (
	"fmt"
	"github.com/adpalmer/go-cluster"
	"math"
)

type Point struct {
	x float64
	y float64
}

func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

func euclideanDist(P1, P2 interface{}) float64 {
	p1, p2 := P1.(Point), P2.(Point)
	return math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
}

func updateCluster(c gocluster.Cluster) interface{} {
	tot := Point{0., 0.}
	for i := 0; i < len(c.Entity); i++ {
		tmp := c.Entity[i].(Point)
		tot.x += tmp.x
		tot.y += tmp.y
	}
	tot.x /= float64(len(c.Entity))
	tot.y /= float64(len(c.Entity))
	return tot
}

func main() {
	pts := []Point{Point{1, 1}, Point{2., .999}, Point{1, .3}, Point{10, 15}, Point{20, 10}, Point{15, 15}}

	// make it work with function
	interfacePoints := make([]interface{}, len(pts))

	for i, d := range pts {
		interfacePoints[i] = d
	}

	// Solve with standard k-means
	km := gocluster.Km(euclideanDist, updateCluster)
	clusterCenters, clusterMembers, _ := km(interfacePoints, 2, 10)
	for i := 0; i < len(clusterCenters); i++ {
		fmt.Printf("Cluster Center %v:\n", clusterCenters[i].(Point))
		fmt.Println("\tmembers -> ", clusterMembers[i].Entity)
	}

	// Solve with k-means++
	km = gocluster.Kmpp(euclideanDist, updateCluster)
	clusterCenters, clusterMembers, _ = km(interfacePoints, 2, 10)
	for i := 0; i < len(clusterCenters); i++ {
		fmt.Printf("Cluster Center %v:\n", clusterCenters[i].(Point))
		fmt.Println("\tmembers -> ", clusterMembers[i].Entity)
	}

}
