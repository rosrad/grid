//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
)

type FeatChk struct {
	Model
}

func NewFeatChk() *FeatChk {
	return &FeatChk{*NewModel()}
}

func (chk FeatChk) Check(set string) error {
	dirs, err := chk.Dst.Subsets(set)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		cmd_str := JoinArgs(
			"feat-to-len",
			"scp:"+path.Join(chk.Src.DataDir(), dir, "feats.scp"))
		// if err := LogRun(cmd_str, g.DecodeDir(dir)); err != nil {
		// 	Err().Println(err)
		// }
		Trace().Println(cmd_str)
	}
	return nil
}

func (chk FeatChk) Decode(set string) error {
	return chk.Check(set)
}

type FeatChkTask struct {
	FeatChk
}

func NewFeatChkTask() *FeatChkTask {
	return &FeatChkTask{*NewFeatChk()}
}

func (t FeatChkTask) Identify() string {
	return "CMVN"
}

func (c FeatChkTask) Run() error {
	return nil
}

func FeatChkTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewFeatChkTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("Cmvn Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
