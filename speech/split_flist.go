package main

//package speech

import (
	"flag"
	"fmt"
)

func main() {
	var head string
	flag.StringVar(&head, "h", "", "append the head for paths ")
	flag.Parse()
	fmt.Println("Test Combine ", head)

}
