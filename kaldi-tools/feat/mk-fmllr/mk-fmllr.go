package main

import (
	"flag"
	"github.com/rosrad/kaldi"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "tri2.cmn", "tri2.cmn(default)")
	flag.PrintDefaults()
	flag.Parse()

	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("mk-fmllr")

	f := kaldi.NewFmllrTask()
	f.Src.Name = name
	f.Src.Exp = "GMM"

	kaldi.WriteTask(f)
	if err := f.Run(); err != nil {
		kaldi.Err().Println()
	}

}
