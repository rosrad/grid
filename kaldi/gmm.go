package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
)

type GmmConf struct {
	Jobs int
	Tri1 bool
}

func NewGmmConf() *GmmConf {
	return &GmmConf{Jobs: 10, Tri1: false}
}

func (conf GmmConf) OptStr() string {
	var_opt := ""
	if conf.Tri1 {
		var_opt = JoinArgs(var_opt, "--boost-silence", "1.25")
	}
	return var_opt
}

type Gmm struct {
	Model
	ModelConf
	GmmConf
}

func NewGmm(feat string) *Gmm {
	g := &Gmm{*NewModel(), *NewModelConf(), *NewGmmConf()}
	// src model values
	g.Src.Feat = feat
	g.Src.Exp = g.Identify()
	g.Src.Label = "normal"
	g.Src.Name = "tri1"
	// copy from src
	g.Dst = g.Src
	g.Dst.Name = "tri2"
	// mk align config from src
	// g.Src = *g.Src.MkAlign()
	return g
}

func (g Gmm) TargetDir() string {
	return g.Dst.ExpDir()
}

func (g Gmm) AlignDir() string {
	return g.Src.AlignDir()
}

func (g Gmm) Subsets(set string) ([]string, error) {
	return g.Dst.Subsets(set)
}

func (g Gmm) DecodeDir(set string) string {
	return path.Join(g.TargetDir(), JoinParams("decode", "#"+path.Base(set)))
}

func (g Gmm) OptStr() string {
	return JoinArgs(g.ModelConf.OptStr(), g.GmmConf.OptStr())
}

func (g Gmm) Train() error {
	gauss_conf := "2500 15000"
	if g.GmmConf.Tri1 {
		gauss_conf = "2000 10000"
	}
	cmd_str := JoinArgs(
		"steps/train_deltas.sh",
		g.OptStr(),
		gauss_conf,
		g.Dst.TrainData("mc"),
		Lang(),
		g.AlignDir(),
		g.TargetDir(),
	)
	err := BashRun(cmd_str)
	if err != nil {
		return err
	}
	return nil
}

// implement the Decoder interface
func (g Gmm) Decode(set string) error {
	dirs, err := g.Dst.Subsets(set)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		cmd_str := JoinArgs(
			"steps/decode.sh",
			"--nj "+strconv.Itoa(JobNum("decode")),
			Graph(g.TargetDir()),
			path.Join(g.Dst.DataDir(), dir),
			g.DecodeDir(dir))
		if err := BashRun(cmd_str); err != nil {
			return err
		}
	}
	return nil
}

func (g Gmm) Score(set string) ([][]string, error) {
	return AutoScore(g.Identify(), DecodeDirs(set, g))
}

func (g Gmm) Identify() string {
	return "GMM"
}

type GmmTask struct {
	Gmm
	TaskBase *TaskBase
	TaskConf *TaskConf
}

func NewGmmTask(feat string) *GmmTask {
	return &GmmTask{*NewGmm(feat), NewTaskBase("gmm", ""), NewTaskConf()}
}

func (t GmmTask) Identify() string {
	return t.Gmm.Identify()
}

func (t GmmTask) Run() error {
	return Run(t.TaskConf, t.Gmm)
}

func GmmTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewGmmTask("mfcc")
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
