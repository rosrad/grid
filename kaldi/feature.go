//
package kaldi

import (
	"path"
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
	basedir   string
}

func NewFeat() *Feat {
	return &Feat{"raw", []string{}, *NewNorm(), 0, true, ""}
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

func (f Feat) Str() string {
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
	return str
}

func (f Feat) TransformStr(dir string) string {
	opt := ""
	for _, mat := range f.Transform {

		if len(strings.Trim(mat, " ")) == 0 {
			continue
		}

		if opt != "" {
			opt += " | "
		}
		opt += JoinArgs("transform-feats",
			path.Join(dir, mat),
			"ark:-", "ark:-")
	}
	return opt
}

func (f Feat) FeatStr(dir string) string {
	str := f.Str()

	if transform := f.TransformStr(dir); transform != "" {
		str = JoinArgs(str, "|", transform)
	}

	return "'" + str + "|'"
}

func (f Feat) FeatOpt(dir string) string {
	return JoinArgs("--feat", f.FeatStr(dir))
}
