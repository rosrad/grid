package main

import (
	"flag"
	"github.com/rosrad/kaldi"
)

func main() {
	var dy string
	var vars bool
	flag.StringVar(&dy, "dy", "raw", "raw(default),delta")
	flag.BoolVar(&vars, "vars", false, "flase(default)")
	flag.PrintDefaults()
	flag.Parse()
	kaldi.Init()
	defer kaldi.Uninit()

	kaldi.Trace().Println("mk-bnf")
	b := kaldi.NewBnfTask()
	b.ModelConf.Dynamic = dy
	b.ModelConf.Norm.Cmvn.Vars = vars
	kaldi.WriteTask(b)
	if err := b.Run(); err != nil {
		kaldi.Err().Println()
	}

}
