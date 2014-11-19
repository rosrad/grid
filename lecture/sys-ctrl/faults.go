//
package main

import (
	"fmt"
	"math"
	"os"
)

func Short(e float64) float64 {
	return 1 - math.Pow(1-math.Pow(e, 2), 3)
}

func Open(e float64) float64 {
	return math.Pow(1-math.Pow(1-e, 2), 3)
}

func main() {
	fmt.Println("faults")
	sf, _ := os.Create("short.dat")
	defer sf.Close()
	of, _ := os.Create("open.dat")
	defer of.Close()

	for i := 0.0; i <= 1.0; i = i + 0.01 {
		fmt.Fprintf(sf, "%.4f, %.4f \n", i, Short(i))
		fmt.Fprintf(of, "%.4f, %.4f \n", i, Open(i))
	}

}
