//
package kaldi

import (
	"encoding/json"
	"fmt"
	"io"
)

type Align struct {
	ExpBase
	Feat
	Fmllr bool
	Extra Extra
}

func NewAlign(feat string) *Align {
	return &Align{*NewExpBase(), *NewFeat(), false, Extra{"", ""}}
}

func (a Align) OptStr() string {
	return JoinArgs(a.Extra.Opts, a.FeatOpt(a.ExpDir()))
}

func (a Align) TrainData() string {
	return a.ExpBase.TrainData(a.Feat.Condition())
}

func (a Align) AlignCmd() (error, string) {
	gmm_align := "steps/align_si.sh"
	if a.Fmllr {
		gmm_align = "steps/align_fmllr.sh"
	}
	switch a.Exp {
	case "MONO":
		return nil, gmm_align
	case "GMM":
		return nil, gmm_align
	case "LDA":
		return nil, gmm_align
	case "SAT":
		return nil, gmm_align
	case "DNN":
		return nil, "steps/nnet2/align.sh"
	case "NET":
		return nil, "steps/nnet/align.sh"
	default:
		return fmt.Errorf("No Effective Align for :%s", a.Exp), ""
	}

}

func (a Align) MkAlign() error {

	err, script := a.AlignCmd()
	if err != nil {
		Err().Println(err)
		return err
	}
	cmd_str := JoinArgs(
		script,
		"--nj", MaxNum(a.TrainData()),
		a.OptStr(),
		a.Extra.Args,
		a.TrainData(),
		Lang(),
		a.ExpDir(),
		a.AlignDir())
	if err := LogCpuRun(cmd_str, a.AlignDir()); err != nil {
		return err
	}
	SyncMat(a.Transform, a.ExpDir(), a.AlignDir())
	return nil
}

type AlignTask struct {
	Align
}

func NewAlignTask(feat string) *AlignTask {
	return &AlignTask{*NewAlign(feat)}
}

func (t AlignTask) Identify() string {
	return "ALIGN"
}

func (t AlignTask) Run() error {
	return t.MkAlign()
}

func AlignTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewAlignTask("mfcc")
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
