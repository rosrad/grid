package main

import (
	"github.com/rosrad/kaldi"
	"net"
	"net/rpc"
)

func main() {
	kaldi.Init("", "")
	defer kaldi.Uninit()
	kaldi.DevInstance().AutoSync()
	kaldi.DevInstance().SortGpu()
	kaldi.DevInstance().PrintNodes(true)
	kaldi.DevInstance().SortCpu()
	kaldi.DevInstance().PrintNodes(false)

	rpc.Register(kaldi.DevInstance().RPC())
	l, e := net.Listen("tcp", ":9001")
	if e != nil {
		kaldi.Err().Println("Grid Svr listen error:", e)
	}

	rpc.Accept(l)
}
