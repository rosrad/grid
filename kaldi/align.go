//
package kaldi

import (
	"encoding/json"
	"io"
	"strconv"
)

type Align struct {
	ExpBase
	Dynamic string
}

func NewAlign(feat string) *Align {
	return &Align{*NewExpBase(), "delta"}
}

func (a Align) OptStr() string {
	return JoinArgs("--feat-type", a.Dynamic)
}

func (a Align) MkAlign() error {
	cmd_str := JoinArgs(
		"steps/align_si.sh",
		"--nj "+strconv.Itoa(JobNum("decode")),
		a.OptStr(),
		a.TrainData("mc"),
		Lang(),
		a.ExpDir(),
		a.AlignDir())
	if err := BashRun(cmd_str); err != nil {
		return err
	}
	return nil
}

type AlignTask struct {
	Align
	TaskBase *TaskBase
}

func NewAlignTask(feat string) *AlignTask {
	return &AlignTask{*NewAlign(feat), NewTaskBase("align", "")}
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
