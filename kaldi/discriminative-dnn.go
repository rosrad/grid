//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
)

// training config struct
type DiscDnnConf struct {
	Criterion string
	Jobs      int
	Gpu       bool
}

func NewDiscDnnConf() *DiscDnnConf {
	return &DiscDnnConf{
		Criterion: "smbr",
		Jobs:      16,
		Gpu:       true}
}

func (conf *DiscDnnConf) OptStr() string {
	var_opt := ""
	if conf.Gpu {
		var_opt = JoinArgs(var_opt, "--num-threads", "1")
	}
	var_opt = JoinArgs(var_opt, "--num-jobs-nnet", strconv.Itoa(conf.Jobs))
	var_opt = JoinArgs(var_opt, "--criterion", conf.Criterion)
	return var_opt
}

type DiscDnn struct {
	Model
	DiscDnnConf
}

func NewDiscDnn() *DiscDnn {
	d := &DiscDnn{*NewModel(), *NewDiscDnnConf()}
	d.Dst = d.Src
	d.Dst.Exp = d.Identify()
	return d
}

func (d DiscDnn) AlignDir() string {
	return d.Src.AlignDir()
}

func (d DiscDnn) TargetDir() string {
	return d.Dst.ExpDir()
}

func (d DiscDnn) DenlatsDir() string {
	return path.Join(d.Src.ExpDir(), "denlats")
}

func (d DiscDnn) Subsets(set string) ([]string, error) {
	return d.Dst.Subsets(set)
}

func (d DiscDnn) DecodeDir(dir string) string {
	return MkDecode(d.TargetDir(), dir)
}

func (d DiscDnn) SrcModel() string {
	return path.Join(d.Src.ExpDir(), "final.mdl")
}

func (d DiscDnn) OptStr() string {
	return JoinArgs(
		d.DiscDnnConf.OptStr(),
		d.Model.OptStr())
}

func (d DiscDnn) DenlatsOptStr() string {
	return JoinArgs(
		d.Model.Feat.OptStr(),
		"--nj", MaxNum(d.Dst.TrainData("mc")),
		"--sub-split 20",
		"--num-threads 6")
}

func (d DiscDnn) MkDenlats() error {
	cmd_str := JoinArgs(
		"steps/nnet2/make_denlats.sh",
		d.DenlatsOptStr(),
		d.Dst.TrainData("mc"),
		Lang(),
		d.Src.ExpDir(),
		d.DenlatsDir())
	Trace().Println(cmd_str)
	return LogGpuRun(cmd_str, d.DenlatsDir())
}

func (d DiscDnn) Train() error {
	if err := d.MkDenlats(); err != nil {
		Err().Println(err)
		return err
	}

	cmd_str := JoinArgs(
		"steps/nnet2/train_discriminative.sh",
		"--learning-rate 0.000002",
		d.OptStr(),
		d.Dst.TrainData("mc"),
		Lang(),
		d.AlignDir(),
		d.DenlatsDir(),
		d.SrcModel(),
		d.TargetDir())
	Trace().Println(cmd_str)
	return LogGpuRun(cmd_str, d.TargetDir())
}

// implement the Decoder interface
func (d DiscDnn) Decode(set string) error {
	items, err := d.Src.Subsets(set)
	if err != nil {
		return err
	}

	for _, dir := range items {
		cmd_str := JoinArgs(
			"steps/nnet2/decode.sh",
			"--nj", MaxNum(path.Join(d.Dst.DataDir(), dir)),
			d.Model.OptStr(),
			Graph(d.TargetDir()),
			path.Join(d.Dst.DataDir(), dir),
			d.DecodeDir(dir))
		if err := LogGpuRun(cmd_str, d.DecodeDir(dir)); err != nil {
			Err().Println(err)
			continue
		}
	}

	return nil
}

// implement the Counter interface
func (d DiscDnn) Score(set string) ([][]string, error) {
	return AutoScore(d.Identify(), DecodeDirs(set, d))
}

func (d DiscDnn) Identify() string {
	return "DISCDNN"
}

type DiscDnnTask struct {
	DiscDnn            //DiscDnn struct that was used
	TaskConf *TaskConf // task config for dnn
}

func NewDiscDnnTask() *DiscDnnTask {
	return &DiscDnnTask{*NewDiscDnn(), NewTaskConf()}
}

func (t DiscDnnTask) Identify() string {
	return t.DiscDnn.Identify()
}

func (t DiscDnnTask) Run() error {
	return Run(t.TaskConf, t.DiscDnn)
}

func DiscDnnTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewDiscDnnTask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
