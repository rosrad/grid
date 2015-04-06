//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"sync"
)

type Plda struct {
	Model
	Ubm ExpBase
}

func NewPlda() *Plda {
	p := Plda{*NewModel(), *NewExpBase()}
	p.Dst.Exp = p.Identify()
	return &p
}

func (p Plda) TargetDir() string {
	return p.Dst.ExpDir()
}

func (p Plda) SourceData() string {
	return p.Src.DataDir()
}

func (p Plda) AlignDir() string {
	return p.Ali.AlignDir()
}

func (p Plda) Subsets(set string) ([]string, error) {
	return p.Src.Subsets(set)
}

func (p Plda) DecodeDir(set string) string {
	return MkDecode(p.TargetDir(), set)
}

func (p Plda) OptStr() string {
	return JoinArgs(p.Feat.OptStr())
}

func (p Plda) MkPreCond() error {
	files := [...]string{"final.ubm", "splice_opts", "cmvn_opts"}
	for _, f := range files {
		InsureDir(p.Dst.ExpDir())
		cmd_str := JoinArgs("cp",
			path.Join(p.Ubm.ExpDir(), f),
			p.Dst.ExpDir()+"/")
		if err := BashRun(cmd_str); err != nil {
			Err().Println(err)
		}
	}

	cmd_str := JoinArgs("cp",
		path.Join(p.Src.ExpDir(), "tree"),
		p.Dst.ExpDir()+"/")
	if err := BashRun(cmd_str); err != nil {
		Err().Println(err)
	}
	return nil
}

func (p Plda) Train() error {

	if err := p.MkPreCond(); err != nil {
		return nil
	}

	cmd_str := JoinArgs(
		"steps/train_plda.sh",
		p.OptStr(),
		"2400 20000",
		p.Src.TrainData(p.Condition()),
		Lang(),
		p.AlignDir(),
		p.TargetDir(),
	)

	if err := LogCpuRun(cmd_str, p.TargetDir()); err != nil {
		return err
	}
	return nil
}

// implement the Decoder interface
func (p Plda) Decode(set string) error {
	dirs, err := p.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		cmd_str := JoinArgs(
			"steps/decode_plda.sh",
			"--stage -1",
			"--nj ", JobNum("decode"),
			p.FeatOpt(),
			Graph(p.TargetDir()),
			path.Join(p.SourceData(), dir),
			p.DecodeDir(dir))
		go func(cmd, dir string) {
			defer wg.Done()
			if err := LogCpuRun(cmd, dir); err != nil {
				Err().Println(err)
			}
		}(cmd_str, p.DecodeDir(dir))
	}
	wg.Wait()
	return nil
}

func (p Plda) Score(set string) ([][]string, error) {
	return AutoScore(p.Identify(), DecodeDirs(set, p))
}

func (p Plda) Identify() string {
	return "PLDA"
}

type PldaTask struct {
	Plda
	TaskConf *TaskConf
}

func NewPldaTask() *PldaTask {
	return &PldaTask{*NewPlda(), NewTaskConf()}
}

func (t PldaTask) Identify() string {
	return t.Plda.Identify()
}

func (t PldaTask) Run() error {
	return Run(t.TaskConf, t.Plda)
}

func PldaTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewPldaTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("PLDA Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
