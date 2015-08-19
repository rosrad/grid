//
package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
)

func main() {
	var update bool
	flag.BoolVar(&update, "u", false, "update server")
	flag.Parse()
	kaldi.Init("", "")
	defer kaldi.Uninit()
	if !kaldi.GridClient().Inited() {
		kaldi.Err().Println("Grid Client no initlized!")
		return
	}
	if update {
		fmt.Println("Update grid nodes")
		kaldi.GridClient().Update()
	}
	fmt.Println("========= Gpu Sort =========")
	for _, n := range kaldi.GridClient().GpuSort() {
		fmt.Printf("Node:%s, GpuUsage:%04.2f, GpuMem:%04.2f \n", n.Node, n.GpuUsage(), n.GpuMem)
	}

	fmt.Println("========= Cpu Sort =========")
	for _, n := range kaldi.GridClient().CpuSort() {
		fmt.Printf("Node:%s, CpuUsage:%04.2f, CPU(s):%02d, CPU MHz: %04.2f, LoadAve:%02.2f\n",
			n.Node, n.CpuUsage(), n.CpuNum, n.CpuMHz, n.LoadAve)
	}
}
