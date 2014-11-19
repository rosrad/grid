//
package main

import (
	"flag"
	"github.com/rosrad/kaldi"
)

func main() {
	var feat, dy, tag string
	flag.StringVar(&feat, "feat", "mfcc", "mfcc, bnf(default)")
	flag.StringVar(&dy, "dy", "delta", "raw, delta(default) ")
	flag.StringVar(&tag, "tag", "", "anything, normal(default)")
	flag.PrintDefaults()
	flag.Parse()

	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("mk-align")
	a := kaldi.NewAlignTask(feat)
	a.Dynamic = dy
	for _, name := range flag.Args() {
		kaldi.Trace().Println("Make alignment of name :", name)
		a.Name = name
		kaldi.WriteTask(a)
		if err := a.Run(); err != nil {
			kaldi.Err().Println()
		}
	}
}
