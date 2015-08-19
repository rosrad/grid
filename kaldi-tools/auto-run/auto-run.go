package main

import (
	"flag"
	"github.com/rosrad/kaldi"
	"os"
	"strings"
)

func main() {
	var manual bool
	var gpu bool
	flag.BoolVar(&manual, "manual", false, "manual to control servers (default=false) ")
	flag.BoolVar(&gpu, "g", false, "gpu dependenc")
	flag.Parse()
	kaldi.Init("", "")
	defer kaldi.Uninit()

	if flag.NArg() == 0 {
		kaldi.Trace().Println("No enough args!")
		return
	}

	cmd := strings.Join(flag.Args(), " ")
	var err error

	if manual {
		err = kaldi.BashRun(cmd)
	} else {
		if gpu {
			err = kaldi.GpuBashRun(cmd)
		} else {
			err = kaldi.CpuBashRun(cmd)
		}
	}

	if err != nil {
		kaldi.Err().Println("Cmd Err:", cmd, "\t", err)
		os.Exit(1)
	}
}
