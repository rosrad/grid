package ml

import (
	"fmt"
)

type Feature struct {
	sum int
	x1  map[string]int
	x2  map[string]int
}

func mean(f *Feature) float64 {
	sum := 0
	var list []int
	for key, count := range f.x1 {
		fmt.Println("Key:", key, "Count:", count)
		list = append(list, count)
		sum += count
	}
	mean := 0.0
	for count := range list {
		p := float64(count) / float64(sum)
		mean += p * p
	}

	return mean
}

type Node struct {
	parent   *Node
	children []*Node
	feature  Feature
}

func test() {
	fmt.Printf("Hello, 世界\n")
}
