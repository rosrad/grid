package kaldi

import (
	"encoding/json"
	"io"
	"strconv"
)

type SplicedGmmConf struct {
	Context int
}

func NewSplicedGmmConf() *SplicedGmmConf {
	return &SplicedGmmConf{Context: 2}
}

func (lc SplicedGmmConf) OptStr() string {
	opt := ""
	if lc.Context > 0 {
		opt = JoinArgs(opt,
			"--splice-opts",
			"\"--left-context="+strconv.Itoa(lc.Context)+" --right-context="+strconv.Itoa(lc.Context)+"\"")
	}
	return opt
}

type SplicedGmm struct {
	Model
	SplicedGmmConf
}

func NewSplicedGmm() *SplicedGmm {
	sg := &SplicedGmm{*NewModel(), *NewSplicedGmmConf()}
	sg.Dst.Exp = sg.Identify()
	return sg
}

func (sg SplicedGmm) Identify() string {
	return "SPLICEDGMM"
}

func (sg SplicedGmm) TargetDir() string {
	return sg.Dst.ExpDir()
}

func (sg SplicedGmm) AlignDir() string {
	return sg.Src.AlignDir()
}

func (sg SplicedGmm) Subsets(set string) ([]string, error) {
	return sg.Dst.Subsets(set)
}

func (sg SplicedGmm) DecodeDir(set string) string {
	return MkDecode(sg.TargetDir(), set)
}

func (sg SplicedGmm) OptStr() string {
	return JoinArgs(sg.Model.OptStr(), sg.SplicedGmmConf.OptStr())
}

func (sg SplicedGmm) Train() error {
	cmd_str := JoinArgs(
		"steps/train_splice.sh",
		sg.OptStr(),
		"1200 10000",
		sg.Dst.TrainData("mc"),
		Lang(),
		sg.AlignDir(),
		sg.TargetDir(),
	)
	err := LogCpuRun(cmd_str, sg.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

type SplicedGmmTask struct {
	SplicedGmm
}

func NewSplicedGmmTask() *SplicedGmmTask {
	return &SplicedGmmTask{*NewSplicedGmm()}
}

func (spt SplicedGmmTask) Identify() string {
	return spt.SplicedGmm.Identify()
}

func (spt SplicedGmmTask) Run() error {
	return spt.Train()
}

func SplicedGmmTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewSplicedGmmTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
