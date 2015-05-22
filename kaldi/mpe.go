//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
	"sync"
)

type MPE struct {
	Model
	NumIter int
}

func NewMPE() *MPE {
	n := &MPE{*NewModel(), 1}
	// src model values
	n.Dst = n.Src
	n.Dst.Exp = n.Identify()
	return n
}

func (n MPE) AlignDir() string {
	return n.Ali.AlignDir()
}

func (n MPE) TargetDir() string {
	return n.Dst.ExpDir()
}

func (n MPE) Subsets(set string) ([]string, error) {
	return n.Src.Subsets(set)
}

func (n MPE) DecodeDir(dir string) string {
	return MkDecode(n.TargetDir(), dir)
}

func (n MPE) DenLatsDir() string {
	return n.Ali.DenLatsDir()
}

func (n MPE) Train() error {
	cmd_str := JoinArgs(
		"steps/nnet/train_mpe.sh",
		"--do-smbr true",
		"--num-iters", strconv.Itoa(n.NumIter),
		n.FeatOpt(n.AlignDir()),
		n.TrainData(),
		Lang(),
		n.Ali.ExpDir(),
		n.AlignDir(),
		n.DenLatsDir(),
		n.TargetDir())
	Trace().Println(cmd_str)
	err := LogGpuRun(cmd_str, n.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

// implement the Decoder interface
func (n MPE) Decode(set string) error {
	items, err := n.Src.Subsets(set)
	if err != nil {
		return err
	}

	decode_opt := ""
	if len(n.DecodeConf) != 0 {
		decode_opt = JoinArgs("--config", n.DecodeConf)
	}

	var wg sync.WaitGroup
	for i := 1; i <= n.NumIter; i++ {
		for _, dir := range items {
			cmd_str := JoinArgs(
				DecodeCmd(n.Identify()),
				decode_opt,
				"--nnet", strconv.Itoa(i)+".nnet",
				"--nj", MaxNum(path.Join(n.Src.DataDir(), dir)),
				n.FeatOpt(n.AlignDir()),
				Graph(n.TargetDir()),
				path.Join(n.Src.DataDir(), dir),
				n.DecodeDir(dir))
			wg.Add(1)
			go func(cmd_str, dir string) {
				defer wg.Done()
				if err := LogCpuRun(cmd_str, dir); err != nil {
					Err().Println(err)
				}
			}(cmd_str, n.DecodeDir(dir))
		}
	}
	wg.Wait()
	return nil
}

// implement the Counter interface
func (n MPE) Score(set string) ([][]string, error) {
	return AutoScore(n.Identify(), DecodeDirs(set, n))
}

func (n MPE) Identify() string {
	return "MPE"
}

type MPETask struct {
	MPE                //MPE struct that was used
	TaskConf *TaskConf // task config for dnn
}

func NewMPETask() *MPETask {
	return &MPETask{*NewMPE(), NewTaskConf()}
}

func (t MPETask) Identify() string {
	return t.MPE.Identify()
}

func (t MPETask) Run() error {
	return Run(t.TaskConf, t.MPE)
}

func MPETasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewMPETask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
