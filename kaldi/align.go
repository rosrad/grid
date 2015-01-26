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
	Extra string
}

func NewAlign(feat string) *Align {
	return &Align{*NewExpBase(), *NewFeat(), ""}
}

func (a Align) OptStr() string {
	return a.FeatOpt()
}

func (a Align) TrainData() string {
	cond := "cln"
	if a.MC {
		cond = "mc"
	}

	return a.ExpBase.TrainData(cond)
}

func (a Align) AlignCmd() (error, string) {

	switch a.Exp {
	case "GMM":
		return nil, "steps/align_si.sh"
	case "DNN":
		return nil, "steps/nnet2/align.sh"
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
		"--nj", JobNum("train"),
		a.OptStr(),
		a.TrainData(),
		Lang(),
		a.ExpDir(),
		a.AlignDir())
	if err := LogCpuRun(cmd_str, a.AlignDir()); err != nil {
		return err
	}

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
