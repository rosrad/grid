package kaldi

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sync"
)

type FBank struct {
	ExpBase
}

func NewFBank() *FBank {
	return &FBank{*NewExpBase()}
}

func (f FBank) Compute(set string) error {
	// copy data scp to feats dir
	InsureDir(f.DataDir())
	cp_cmd := JoinArgs(
		"rsync -avz",
		path.Join("data", set),
		f.DataDir()+"/")
	fmt.Println(cp_cmd)
	BashRun(cp_cmd)
	cmd := JoinArgs(
		"steps/make_fbank.sh",
		"--nj", MaxNum(path.Join(f.DataDir(), set)),
		path.Join(f.DataDir(), set),
		path.Join(f.LogDir(), "make_fbank", set+".log"),
		path.Join(f.ParamDir(), set))
	// fmt.Println(cmd)
	LogCpuRun(cmd, f.DataDir())
	return nil
}

func (f FBank) ComputeAll(set DataType) {
	var wg sync.WaitGroup
	for _, set := range DataSets(set) {
		wg.Add(1)
		go func(set string) {
			defer wg.Done()
			f.Compute(set)
		}(set)
	}
	wg.Wait()
}

type FBankTask struct {
	FBank
}

func NewFBankTask() *FBankTask {
	return &FBankTask{*NewFBank()}
}

func (t FBankTask) Identify() string {
	return "MFBANK"
}

func (t FBankTask) Run() error {
	set := TRAIN_MC_TEST
	if len(SysConf().DecodeSet) != 0 {
		set = Str2DataType(SysConf().DecodeSet)
	}
	t.ComputeAll(set)
	return nil
}

func FBankTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewFBankTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("FBank Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
