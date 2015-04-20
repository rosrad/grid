package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"sync"
)

type Gmm struct {
	Model
}

func NewGmm() *Gmm {
	g := &Gmm{*NewModel()}
	// src model values
	g.Src.Exp = g.Identify()
	g.Dst = g.Src
	return g
}

func (g Gmm) TargetDir() string {
	return g.Dst.ExpDir()
}

func (g Gmm) AlignDir() string {
	return g.Ali.AlignDir()
}

func (g Gmm) Subsets(set string) ([]string, error) {
	return g.Src.Subsets(set)
}

func (g Gmm) DecodeDir(set string) string {
	return MkDecode(g.TargetDir(), set)
}

func (g Gmm) OptStr() string {
	var_opt := ""
	if !g.MC {
		var_opt = JoinArgs(var_opt, "--boost-silence", "1.25")
	}

	return JoinArgs(var_opt, g.Feat.OptStr())
}

func (g Gmm) Gaussian() string {
	gaussian := "2000 10000"
	if g.MC {
		gaussian = "2500 15000"
	}
	return gaussian
}

func (g Gmm) Train() error {
	cmd_str := JoinArgs(
		"steps/train_deltas.sh",
		g.OptStr(),
		g.Gaussian(),
		g.Src.TrainData(g.Condition()),
		Lang(),
		g.AlignDir(),
		g.TargetDir(),
	)

	err := LogCpuRun(cmd_str, g.TargetDir())
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
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		cmd_str := JoinArgs(
			"steps/decode.sh",
			"--nj ", MaxNum(path.Join(g.Src.DataDir(), dir)),
			g.FeatOpt(),
			Graph(g.TargetDir()),
			path.Join(g.Src.DataDir(), dir),
			g.DecodeDir(dir))
		go func(cmd, dir string) {
			defer wg.Done()
			if err := LogCpuRun(cmd, dir); err != nil {
				Err().Println(err)
			}
		}(cmd_str, g.DecodeDir(dir))
	}
	wg.Wait()
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
	TaskConf *TaskConf
}

func NewGmmTask() *GmmTask {
	return &GmmTask{*NewGmm(), NewTaskConf()}
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
		t := NewGmmTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
