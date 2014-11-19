package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
)

func main() {
	fmt.Println("cmvn-comp")
	var mode, feat string
	flag.StringVar(&mode, "mode", "utt", "mode name:fake, utt(per-utterance), spk(per-speaker)")
	flag.StringVar(&feat, "feat", "mfcc", "feature type:mfcc or bnf")

	flag.Parse()
	dataset := flag.Args()
	if len(dataset) == 0 {
		dataset = kaldi.DataSets(false)
		fmt.Println("All DataSets were used here!")
	}

	fmt.Println("Feat:", feat)
	fmt.Println("Mode:", mode)
	fmt.Println("DataSets:", dataset)
	opt := kaldi.NewCmvnOption(feat, mode)
	fmt.Println("opt:", opt)

	for _, set := range dataset {
		subsets, _ := kaldi.Subsets(set, feat)
		for _, subset := range subsets {
			opt.CompCmvn(subset)
		}
	}
	fmt.Println("cmvn-comp finished!")
}
