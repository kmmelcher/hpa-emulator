package main

import (
	"fmt"

	"github.com/softchris/math"
	podautoscaler "k8s.io/kubernetes/pkg/controller/podautoscaler"
)

func main(){
	fmt.Println("oi")

	x := math.Add(1,1)

	fmt.Println(x)

	podautoscaler.NewHorizontalController()

}
