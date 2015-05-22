//
package kaldi

import (
	"path"
)

type Extra struct {
	Opts string
	Args string
}

type Model struct {
	Dst   ExpBase
	Src   ExpBase
	Ali   ExpBase
	Extra Extra
	Feat
	DecodeConf string
}

func (m Model) TrainData() string {
	return m.Src.TrainData(m.Feat.Condition())
}

func (m Model) OptStr() string {
	opt := JoinArgs(
		m.Extra.Opts,
		"--nj", MaxNum(m.TrainData()),
		m.FeatOpt(m.Ali.ExpDir()))
	return opt
}
func NewModel() *Model {
	return &Model{*NewExpBase(), *NewExpBase(), *NewExpBase(), Extra{"", ""}, *NewFeat(), ""}
}

func SyncMat(transform []string, src, dst string) {
	if len(transform) == 0 {
		return
	}
	for _, f := range transform {
		cmd_str := JoinArgs("rsync",
			"-avzLr",
			path.Join(src, f),
			dst+"/")
		Trace().Println(cmd_str)
		err := CpuBashRun(cmd_str)
		if err != nil {
			Err().Println(err)
		}
	}
}
