//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
)

// training config struct
type DnnConf struct {
	MinBatch int
	Jobs     int
	Layers   int
	Context  int
	Gpu      bool
}

func NewDnnConf() *DnnConf {
	return &DnnConf{
		MinBatch: 512,
		Jobs:     JobNum("dnn"),
		Layers:   2,
		Context:  1,
		Gpu:      true}
}

func (conf *DnnConf) OptStr() string {
	var_opt := ""
	if conf.Gpu {
		var_opt = JoinArgs(var_opt, "--num-threads", "1")
	}
	if conf.Context != 1 {
		var_opt = JoinArgs(var_opt, "--splice-width", strconv.Itoa(conf.Context))
	}
	var_opt = JoinArgs(var_opt, "--minibatch-size", strconv.Itoa(conf.MinBatch))
	var_opt = JoinArgs(var_opt, "--num-hidden-layers", strconv.Itoa(conf.Layers))
	var_opt = JoinArgs(var_opt, "--num-jobs-nnet", strconv.Itoa(conf.Jobs))
	return var_opt
}

type Dnn struct {
	Model
	ModelConf
	DnnConf DnnConf
}

func NewDnn(feat string) *Dnn {
	d := &Dnn{*NewModel(), *NewModelConf(), *NewDnnConf()}
	// src model values
	d.Src.Feat = feat
	d.Src.Exp = "GMM"
	d.Src.Label = "normal"
	d.Src.Name = "tri1"
	// copy from src
	d.Dst = d.Src
	d.Dst.Exp = d.Identify()
	// mk align config from src
	// d.Src = *d.Src.MkAlign()
	return d
}

func (d Dnn) AlignDir() string {
	return d.Src.AlignDir()
}

func (d Dnn) TargetDir() string {
	return d.Dst.ExpDir()
}

func (d Dnn) Subsets(set string) ([]string, error) {
	return d.Dst.Subsets(set)
}

func (d Dnn) DecodeDir(dir string) string {
	return path.Join(d.TargetDir(), JoinParams("decode#", path.Base(dir)))
}

func (d Dnn) OptStr() string {
	return JoinArgs(
		d.DnnConf.OptStr(),
		d.ModelConf.OptStr())
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
		d.Dst.TrainData("mc"),
		Lang(),
		d.AlignDir(),
		d.TargetDir())
	Trace().Println(cmd_str)

	err := BashRun(cmd_str)
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

	for _, dir := range items {
		cmd_str := JoinArgs(
			"steps/nnet2/decode.sh",
			"--num-threads 1",
			"--nj "+strconv.Itoa(JobNum("decode")),
			Graph(d.TargetDir()),
			path.Join(d.Dst.DataDir(), dir),
			d.DecodeDir(dir))

		if err := BashRun(cmd_str); err != nil {
			return err
		}
	}

	return nil
}

// func (d Dnn) DecodeDirs(set string) []string {
// 	items, err := d.Subsets(set)
// 	dirs := []string{}
// 	if err != nil {
// 		Err().Println("Generate Subset Error:", err)
// 		return dirs
// 	}
// 	for _, item := range items {
// 		dir := d.DecodeDir(item)
// 		if !DirExist(dir) {
// 			continue
// 		}
// 		dirs = append(dirs, dir)
// 	}
// 	return dirs
// }

// implement the Counter interface
func (d Dnn) Score(set string) ([][]string, error) {
	return AutoScore(d.Identify(), DecodeDirs(set, d))
}

func (d Dnn) Identify() string {
	return "DNN"
}

type DnnTask struct {
	Dnn      //Dnn struct that was used
	TaskBase *TaskBase
	TaskConf *TaskConf // task config for dnn
}

func NewDnnTask(feat string) *DnnTask {
	return &DnnTask{*NewDnn(feat), NewTaskBase("dnn", ""), NewTaskConf()}
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
		t := NewDnnTask("mfcc")
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
