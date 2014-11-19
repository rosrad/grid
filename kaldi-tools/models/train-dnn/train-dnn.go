//
package main

import (
	"flag"
	"github.com/rosrad/kaldi"
)

func mkTag(prefix string) string {
	if prefix == "" {
		return kaldi.Now()
	}
	return kaldi.JoinParams(prefix, kaldi.Now())
}

func main() {
	var vars bool
	var feat, dy, tag string
	flag.StringVar(&feat, "feat", "mfcc", "mfcc, bnf(default)")
	flag.StringVar(&dy, "dy", "delta", "raw, delta(default) ")
	flag.StringVar(&tag, "tag", "", "anything, normal(default)")
	flag.BoolVar(&vars, "vars", false, "true, false(default)")
	flag.PrintDefaults()
	flag.Parse()

	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("train-dnn")

	dt := kaldi.NewDnnTask(feat)
	dt.TaskConf.Btrain = true
	dt.TaskConf.Bdecode = true
	dt.Dnn.ModelConf.Dynamic = dy
	dt.Dnn.ModelConf.Norm.Cmvn.Mean = true
	dt.Dnn.ModelConf.Norm.Cmvn.Vars = vars
	kaldi.WriteTask(dt)
	if err := dt.Run(); err != nil {
		kaldi.Err().Println()
	}

}
