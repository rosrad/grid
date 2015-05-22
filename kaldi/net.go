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
type NetConf struct {
	MinBatch   int
	Jobs       int
	Nodes      int
	Context    int
	Layers     int
	Gpu        bool
	NoPretrain bool
}

func NewNetConf() *NetConf {
	return &NetConf{
		Nodes:      500,
		Context:    0,
		Layers:     2,
		Gpu:        true,
		NoPretrain: false}
}

func (conf *NetConf) OptStr() string {
	var_opt := ""
	var_opt = JoinArgs(var_opt, "--splice", strconv.Itoa(conf.Context))
	var_opt = JoinArgs(var_opt, "--nn-depth", strconv.Itoa(conf.Layers))
	var_opt = JoinArgs(var_opt, "--hid-dim", strconv.Itoa(conf.Nodes))
	return var_opt
}

type Net struct {
	Model
	NetConf NetConf
}

func NewNet() *Net {
	n := &Net{*NewModel(), *NewNetConf()}
	// src model values
	n.Dst = n.Src
	n.Dst.Exp = n.Identify()
	return n
}

func (n Net) AlignDir() string {
	return n.Ali.AlignDir()
}

func (n Net) TargetDir() string {
	return n.Dst.ExpDir()
}

func (n Net) PreTargetDir() string {
	return path.Join(n.Dst.ExpDir(), "RBM")
}

func (n Net) Subsets(set string) ([]string, error) {
	return n.Src.Subsets(set)
}

func (n Net) TrainSets() (string, string) {
	tr := path.Join(n.TargetDir(), "tr90")
	cv := path.Join(n.TargetDir(), "cv10")
	return tr, cv
}

func (n Net) DecodeDir(dir string) string {
	return MkDecode(n.TargetDir(), dir)
}

func (n Net) OptStr() string {
	return JoinArgs(
		n.NetConf.OptStr(),
		n.Model.OptStr())
}

func (n Net) SubTrainData() error {
	tr, cv := n.TrainSets()
	cmd_str := JoinArgs(
		"utils/subset_data_dir_tr_cv.sh",
		n.Src.TrainData(n.Condition()),
		tr, cv)
	Trace().Println(cmd_str)
	err := LogCpuRun(cmd_str, n.TargetDir())
	if err != nil {
		return err
	}
	return nil
}
func (n Net) TrainRBM() error {
	cmd_str := JoinArgs(
		"steps/nnet/pretrain_dbn.sh",
		"--rbm-iter 3",
		n.OptStr(),
		n.Src.TrainData(n.Condition()),
		n.PreTargetDir())

	Trace().Println(cmd_str)
	err := LogGpuRun(cmd_str, n.PreTargetDir())
	if err != nil {
		return err
	}
	return nil

}

func (n Net) TrainNet() error {
	dbn := path.Join(n.PreTargetDir(), "final.dbn")
	ft := path.Join(n.PreTargetDir(), "final.feature_transform")
	tr, cv := n.TrainSets()
	cmd_str := JoinArgs(
		"steps/nnet/train.sh",
		"--dbn", dbn,
		"--learn-rate 0.008",
		"--hid-layers 0",
		"--feature-transform", ft,
		n.FeatOpt(n.AlignDir()),
		tr, cv,
		Lang(),
		n.AlignDir(),
		n.AlignDir(),
		n.TargetDir())
	Trace().Println(cmd_str)
	err := LogGpuRun(cmd_str, n.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

func (n Net) Train() error {
	n.SubTrainData()
	if !n.NetConf.NoPretrain {
		err := n.TrainRBM()
		if err != nil {
			return err
		}
	}
	return n.TrainNet()
}

// implement the Decoder interface
func (n Net) Decode(set string) error {
	items, err := n.Src.Subsets(set)
	if err != nil {
		return err
	}
	ft := path.Join(n.PreTargetDir(), "final.feature_transform")

	decode_opt := ""
	if len(n.DecodeConf) != 0 {
		decode_opt = JoinArgs("--config", n.DecodeConf)
	}

	var wg sync.WaitGroup
	for _, dir := range items {
		cmd_str := JoinArgs(
			DecodeCmd(n.Identify()),
			decode_opt,
			"--feature-transform", ft,
			"--nj", MaxNum(path.Join(n.Src.DataDir(), dir)),
			n.FeatOpt(n.AlignDir()),
			Graph(n.TargetDir()),
			path.Join(n.Src.DataDir(), dir),
			n.DecodeDir(dir))
		wg.Add(1)
		go func(cmd_str, dir string) {
			defer wg.Done()
			if err := LogGpuRun(cmd_str, dir); err != nil {
				Err().Println(err)
			}
		}(cmd_str, n.DecodeDir(dir))
	}
	wg.Wait()
	return nil
}

// implement the Counter interface
func (n Net) Score(set string) ([][]string, error) {
	return AutoScore(n.Identify(), DecodeDirs(set, n))
}

func (n Net) Identify() string {
	return "NET"
}

type NetTask struct {
	Net                //Net struct that was used
	TaskConf *TaskConf // task config for dnn
}

func NewNetTask() *NetTask {
	return &NetTask{*NewNet(), NewTaskConf()}
}

func (t NetTask) Identify() string {
	return t.Net.Identify()
}

func (t NetTask) Run() error {
	return Run(t.TaskConf, t.Net)
}

func NetTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewNetTask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
