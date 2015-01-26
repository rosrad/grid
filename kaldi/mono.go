//
package kaldi

import (
	"encoding/json"
	"io"
)

type Mono struct {
	Model
}

func NewMono() *Mono {
	m := &Mono{*NewModel()}
	// src model values
	m.Src.Feat = "mfcc"
	m.Src.Exp = "GMM"
	m.Src.Label = "normal"
	m.Src.Name = "mono"
	m.Dst = m.Src
	m.Dynamic = "delta"
	return m
}

func (m Mono) Identify() string {
	return "MONO"
}

func (m Mono) OptStr() string {
	opt := JoinArgs(
		"--boost-silence 1.25",
		"--nj", JobNum("train"),
		m.Model.OptStr())
	return opt
}

func (m Mono) TargetDir() string {
	return m.Dst.ExpDir()
}

func (m Mono) TrainData() string {
	cond := "cln"
	if m.MC {
		cond = "mc"
	}
	return m.Dst.TrainData(cond)
}

func (m Mono) Train() error {
	cmd_str := JoinArgs(
		"steps/train_mono.sh",
		m.OptStr(),
		m.TrainData(),
		Lang(),
		m.TargetDir())
	Trace().Println(cmd_str)
	err := LogCpuRun(cmd_str, m.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

func (m Mono) MkGraph() error {
	return MkGraphOpt(m.TargetDir(), "--mono")
}

type MonoTask struct {
	Mono
	*TaskConf
}

func NewMonoTask() *MonoTask {
	return &MonoTask{*NewMono(), NewTaskConf()}
}

func (mt MonoTask) Identify() string {
	return mt.Mono.Identify()
}

func (mt MonoTask) Run() error {
	if mt.Btrain {
		if err := mt.Train(); err != nil {
			return err
		}
		if err := mt.MkGraph(); err != nil {
			return err
		}
	}
	return nil
}

func MonoTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewMonoTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
