//
package kaldi

import (
	"encoding/json"
	"io"
)

type DenLats struct {
	ExpBase
	Feat
	Extra Extra
}

func NewDenLats() *DenLats {
	return &DenLats{*NewExpBase(), *NewFeat(), Extra{"", ""}}
}

func (a DenLats) OptStr() string {
	return JoinArgs(a.Extra.Opts, a.FeatOpt(a.ExpDir()))
}

func (a DenLats) TrainData() string {
	return a.ExpBase.TrainData(a.Feat.Condition())
}

func (a DenLats) MkDenLats() error {
	cmd_str := JoinArgs(
		"steps/net/make_denlats.sh",
		"--nj", MaxNum(a.TrainData()),
		a.OptStr(),
		a.Extra.Args,
		a.TrainData(),
		Lang(),
		a.ExpDir(),
		a.DenLatsDir())
	if err := LogCpuRun(cmd_str, a.DenLatsDir()); err != nil {
		return err
	}
	return nil
}

type DenLatsTask struct {
	DenLats
}

func NewDenLatsTask() *DenLatsTask {
	return &DenLatsTask{*NewDenLats()}
}

func (t DenLatsTask) Identify() string {
	return "DENLATS"
}

func (t DenLatsTask) Run() error {
	return t.MkDenLats()
}

func DenLatsTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewDenLatsTask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
