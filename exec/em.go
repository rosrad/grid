package main

import (
	"fmt"
	"math"
	"strconv"
	"util"
)

type Gauss struct {
	X float64
	U float64
	A float64
	W float64
}

func Gaussian(g Gauss) float64 {
	fmt.Println("GS=", g)
	n1 := 1 / math.Sqrt(2*math.Pi*g.A)
	fmt.Println("n1", n1)
	res := n1 * math.Exp(-(g.X-g.U)*(g.X-g.U)/(2*g.A))
	fmt.Println("Gaussian", res)
	return res
}

func Estep(gs *[2]Gauss, data float64) {

}
func main() {

	var gs [2]Gauss
	gs[0].W = 0.5
	gs[0].W = 0.5
	gs[0].U = 1
	gs[1].U = 1
	gs[0].A = 1
	gs[1].A = 0
	lines, _ := util.ReadLines("D:/trainingdata/data2.csv")

	var sum_r1, sum_r2 float64
	var sum_u1, sum_u2 float64 //expection
	var sum_a1, sum_a2 float64 //variance
	for i := 0; i < 10; i++ {
		sum_r1 = 0
		sum_r2 = 0
		sum_u1 = 0
		sum_u2 = 0
		sum_a1 = 0
		sum_a2 = 0
		for _, content := range lines {
			fmt.Println(content)
			x, _ := strconv.ParseFloat(content, len(content))
			fmt.Println("X=", x)
			gs[0].X = x
			gs[1].X = x
			r1 := gs[0].W * Gaussian(gs[0]) / (gs[0].W*Gaussian(gs[0]) + gs[1].W*Gaussian(gs[1]))
			fmt.Println("r1=", r1)
			r2 := 1 - r1
			sum_r1 += r1
			sum_u1 += r1 * x
			sum_a1 += r1 * x * x
			sum_r2 += r2
			sum_u2 += r2 * x
			sum_a2 += r2 * x * x
		}
		fmt.Println(sum_r1, sum_r2)
		gs[0].U = sum_u1 / sum_r1
		gs[0].A = sum_a1 / sum_r1
		gs[0].W = sum_r1 / float64(len(lines))

		gs[1].U = sum_u2 / sum_r2
		gs[1].A = sum_a2 / sum_r2
		gs[1].W = sum_r2 / float64(len(lines))
	}
	for key, item := range gs {
		fmt.Printf("Gaussion%d W = %f        U = %f     A=%f/n", key, item.W, item.U, item.A)

	}

}

