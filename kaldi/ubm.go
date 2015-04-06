//
package kaldi

import (
	"encoding/json"
	"io"
)

type Ubm struct {
	Model
}

func NewUbm() *Ubm {
	u := &Ubm{*NewModel()}
	u.Dst.Exp = u.Identify()
	return u
}

func (u Ubm) Identify() string {
	return "UBM"
}

func (u Ubm) AlignDir() string {
	return u.Ali.ExpDir()
}

func (u Ubm) SourceData() string {
	return u.Src.DataDir()
}

func (u Ubm) TargetDir() string {
	return u.Dst.ExpDir()
}

func (u Ubm) OptStr() string {
	return u.Model.OptStr()
}

func (u Ubm) Train() error {
	cmd_str := JoinArgs(
		"steps/train_ubm_splice.sh",
		u.OptStr(),
		"100",
		u.Src.TrainData("mc"),
		Lang(),
		u.AlignDir(),
		u.TargetDir())
	err := LogCpuRun(cmd_str, u.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

type UbmTask struct {
	Ubm
}

func NewUbmTask() *UbmTask {
	return &UbmTask{*NewUbm()}
}

func (t UbmTask) Identify() string {
	return t.Ubm.Identify()
}

func (t UbmTask) Run() error {
	return t.Train()
}

func UbmTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewUbmTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
