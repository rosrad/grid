//
package main

import (
	"flag"
	"github.com/rosrad/kaldi"
)

func main() {
	var vars bool
	var feat, dy string
	flag.StringVar(&feat, "feat", "mfcc", "mfcc, bnf(default)")
	flag.StringVar(&dy, "dy", "delta", "raw, delta(default) ")
	flag.BoolVar(&vars, "vars", false, "true, false(default)")
	flag.PrintDefaults()
	flag.Parse()

	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("train-gmm")
	g := kaldi.NewGmmTask(feat)
	g.ModelConf.Dynamic = dy
	g.ModelConf.Norm.Cmvn.Mean = true
	g.ModelConf.Norm.Cmvn.Vars = vars
	kaldi.WriteTask(g)
	if err := g.Run(); err != nil {
		kaldi.Err().Println()
	}
}
