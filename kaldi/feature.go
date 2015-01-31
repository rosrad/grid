//
package kaldi

import (
	"strconv"
	"strings"
)

func JobStr() string {
	return "SDATA_JOB"
}

type Feat struct {
	Dynamic   string   // raw,delta
	Transform []string // transform dirs
	Norm      Norm     // Normalision config
	Context   int
	MC        bool
}

func NewFeat() *Feat {
	return &Feat{"raw", []string{}, *NewNorm(), 0, true}
}

func (f Feat) Condition() string {
	cond := "mc"
	if !f.MC {
		cond = "cln"
	}
	return cond
}

func (f Feat) DynamicStr() string {
	if f.Dynamic == "delta" {
		return JoinArgs("add-deltas",
			"ark:-",
			"ark:- ")
	}
	return ""
}

func (f Feat) SpliceStr() string {
	n := strconv.Itoa(f.Context)
	if f.Context > 0 {
		return JoinArgs("splice-feats",
			"--left-context="+n,
			"--right-context="+n,
			"ark:-", "ark:-")
	}
	return ""
}

func (f Feat) FeatStr() string {
	str := JoinArgs("ark,s,cs:copy-feats", "scp:"+JobStr()+"/feats.scp", "ark:-")
	if norm := f.Norm.NormStr(); norm != "" {
		str = JoinArgs(str, "|", norm)
	}
	if splice := f.SpliceStr(); splice != "" {
		str = JoinArgs(str, "|", splice)
	}

	if dy := f.DynamicStr(); dy != "" {
		str = JoinArgs(str, "|", dy)
	}

	if transform := f.TransformStr(); transform != "" {
		str = JoinArgs(str, "|", transform)
	}

	return "'" + str + "|'"
}

func (f Feat) TransformStr() string {
	opt := ""
	for _, dir := range f.Transform {

		if len(strings.Trim(dir, " ")) == 0 {
			continue
		}

		if opt != "" {
			opt += " | "
		}
		opt += JoinArgs("transform-feats",
			dir,
			"ark:-", "ark:-")
	}
	return opt
}

func (f Feat) FeatOpt() string {
	return JoinArgs("--feat", f.FeatStr())
}

func (f Feat) OptStr() string {
	return f.FeatOpt()
}
