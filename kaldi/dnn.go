//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
	"sync"
)

// training config struct
type DnnConf struct {
	MinBatch int
	Jobs     int
	Nodes    int
	Context  int
	Layers   int
	Gpu      bool
}

func NewDnnConf() *DnnConf {
	return &DnnConf{
		MinBatch: 256,
		Jobs:     16,
		Nodes:    500,
		Context:  0,
		Layers:   2,
		Gpu:      true}
}

func (conf *DnnConf) OptStr() string {
	var_opt := ""
	if conf.Gpu {
		var_opt = JoinArgs(var_opt, "--num-threads", "1")
	}
	var_opt = JoinArgs(var_opt, "--minibatch-size", strconv.Itoa(conf.MinBatch))
	var_opt = JoinArgs(var_opt, "--splice-width", strconv.Itoa(conf.Context))
	var_opt = JoinArgs(var_opt, "--num-hidden-layers", strconv.Itoa(conf.Layers))
	var_opt = JoinArgs(var_opt, "--num-jobs-nnet", strconv.Itoa(conf.Jobs))
	var_opt = JoinArgs(var_opt, "--hidden-layer-dim", strconv.Itoa(conf.Nodes))
	return var_opt
}

type Dnn struct {
	Model
	DnnConf DnnConf
}

func NewDnn() *Dnn {
	d := &Dnn{*NewModel(), *NewDnnConf()}
	// src model values
	d.Dst = d.Src
	d.Dst.Exp = d.Identify()
	return d
}

func (d Dnn) AlignDir() string {
	return d.Ali.AlignDir()
}

func (d Dnn) TargetDir() string {
	return d.Dst.ExpDir()
}

func (d Dnn) Subsets(set string) ([]string, error) {
	return d.Src.Subsets(set)
}

func (d Dnn) DecodeDir(dir string) string {
	return MkDecode(d.TargetDir(), dir)
}

func (d Dnn) OptStr() string {
	return JoinArgs(
		d.DnnConf.OptStr(),
		d.Model.OptStr())
}

func (d Dnn) Train() error {
	cmd_str := JoinArgs(
		"steps/nnet2/train_tanh.sh",
		"--mix-up 5000",
		"--initial-learning-rate 0.015",
		"--final-learning-rate 0.002",
		"--num_epochs 20",
		"--num-epochs-extra 10",
		"--add-layers-period 1",
		"--shrink-interval 3",
		d.OptStr(),
		d.Src.TrainData(d.Condition()),
		Lang(),
		d.AlignDir(),
		d.TargetDir())
	Trace().Println(cmd_str)
	err := LogGpuRun(cmd_str, d.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

// implement the Decoder interface
func (d Dnn) Decode(set string) error {
	items, err := d.Src.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, dir := range items {
		cmd_str := JoinArgs(
			DecodeCmd(d.Identify()),
			"--num-threads 1",
			"--nj", MaxNum(path.Join(d.Src.DataDir(), dir)),
			d.FeatOpt(d.AlignDir()),
			Graph(d.TargetDir()),
			path.Join(d.Src.DataDir(), dir),
			d.DecodeDir(dir))
		wg.Add(1)
		go func(cmd_str, dir string) {
			defer wg.Done()
			if err := LogCpuRun(cmd_str, dir); err != nil {
				Err().Println(err)
			}
		}(cmd_str, d.DecodeDir(dir))
	}
	wg.Wait()
	return nil
}

// implement the Counter interface
func (d Dnn) Score(set string) ([][]string, error) {
	return AutoScore(d.Identify(), DecodeDirs(set, d))
}

func (d Dnn) Identify() string {
	return "DNN"
}

type DnnTask struct {
	Dnn                //Dnn struct that was used
	TaskConf *TaskConf // task config for dnn
}

func NewDnnTask() *DnnTask {
	return &DnnTask{*NewDnn(), NewTaskConf()}
}

func (t DnnTask) Identify() string {
	return t.Dnn.Identify()
}

func (t DnnTask) Run() error {
	return Run(t.TaskConf, t.Dnn)
}

func DnnTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewDnnTask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
