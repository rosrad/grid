package main

import (
	"flag"
	"github.com/rosrad/kaldi"
)

func main() {
	var tag string
	flag.StringVar(&tag, "tag", "normal", "normal(default)")
	flag.PrintDefaults()
	flag.Parse()
	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("mk-cmvn")
	if flag.NArg() < 1 {
		kaldi.Err().Println("no enough args for feature type")
	}

	for _, feat := range flag.Args() {
		c := kaldi.NewCmvnTask()
		c.Feat = feat
		kaldi.WriteTask(c)
		if err := c.Run(); err != nil {
			kaldi.Err().Println()
		}
	}

}
