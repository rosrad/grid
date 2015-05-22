package kaldi

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sync"
)

type Mfcc struct {
	ExpBase
}

func NewMfcc() *Mfcc {
	return &Mfcc{*NewExpBase()}
}

func (c Mfcc) Compute(set string) error {
	// copy data scp to feats dir
	InsureDir(c.DataDir())
	cp_cmd := JoinArgs(
		"rsync -avz",
		path.Join("data", set),
		c.DataDir()+"/")
	fmt.Println(cp_cmd)
	BashRun(cp_cmd)
	cmd := JoinArgs(
		"steps/make_mfcc.sh",
		"--nj", MaxNum(path.Join(c.DataDir(), set)),
		path.Join(c.DataDir(), set),
		path.Join(c.LogDir(), "make_mfcc", set+".log"),
		path.Join(c.ParamDir(), set))
	// fmt.Println(cmd)
	LogCpuRun(cmd, c.DataDir())
	return nil
}

func (c Mfcc) ComputeAll(set DataType) {
	Trace().Println(set)
	var wg sync.WaitGroup
	for _, set := range DataSets(set) {
		wg.Add(1)
		go func(set string) {
			defer wg.Done()
			c.Compute(set)
		}(set)
	}
	wg.Wait()
}

type MfccTask struct {
	Mfcc
}

func NewMfccTask() *MfccTask {
	return &MfccTask{*NewMfcc()}
}

func (t MfccTask) Identify() string {
	return "MFCC"
}

func (c MfccTask) Run() error {
	set := TRAIN_MC_TEST
	if len(SysConf().DecodeSet) != 0 {
		set = Str2DataType(SysConf().DecodeSet)
	}
	c.ComputeAll(set)
	return nil
}

func MfccTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewMfccTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("Mfcc Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
