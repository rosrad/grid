//
package main

import (
	"fmt"
	"github.com/rosrad/test/testlib"
)

func TestGmail() {
	fmt.Println(testlib.SendGmail("test go", "from nagaoka"))
}
func TestConcurrency() {
	c := make(chan int)

	// Make the writing operation be performed in
	// another goroutine.
	// go func() {
	// 	c <- 42
	// }()
	c <- 32
	val := <-c
	println(val)
}
func main() {
	TestConcurrency()
}
