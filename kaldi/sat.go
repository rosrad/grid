package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"sync"
)

type Sat struct {
	Model
}

func NewSat() *Sat {
	s := &Sat{*NewModel()}
	// src model values
	s.Extra.Args = "2500 15000"
	s.Src.Exp = s.Identify()
	s.Dst = s.Src
	return s
}

func (s Sat) TargetDir() string {
	return s.Dst.ExpDir()
}

func (s Sat) AlignDir() string {
	return s.Ali.AlignDir()
}

func (s Sat) Subsets(set string) ([]string, error) {
	return s.Src.Subsets(set)
}

func (s Sat) DecodeDir(set string) string {
	return MkDecode(s.TargetDir(), set)
}

func (s Sat) SyncMat() {
	if len(s.Model.Transform) == 0 {
		return
	}
	for _, f := range s.Feat.Transform {
		src := s.Src.ExpDir()
		dst := s.TargetDir()
		cmd_str := JoinArgs("rsync",
			"-avzLr",
			path.Join(src, f),
			dst+"/")
		Trace().Println(cmd_str)
		err := CpuBashRun(cmd_str)
		if err != nil {
			Err().Println(err)
		}
	}
}

func (s Sat) Train() error {
	cmd_str := JoinArgs(
		"steps/train_sat.sh",
		s.OptStr(),
		s.Extra.Args,
		s.TrainData(),
		Lang(),
		s.AlignDir(),
		s.TargetDir(),
	)

	err := LogCpuRun(cmd_str, s.TargetDir())
	if err != nil {
		return err
	}
	SyncMat(s.Model.Transform, s.AlignDir(), s.TargetDir())
	return nil
}

// implement the Decoder interface
func (s Sat) Decode(set string) error {
	dirs, err := s.Dst.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		cmd_str := JoinArgs(
			DecodeCmd(s.Identify()),
			"--nj ", MaxNum(path.Join(s.Src.DataDir(), dir)),
			s.FeatOpt(s.AlignDir()),
			Graph(s.TargetDir()),
			path.Join(s.Src.DataDir(), dir),
			s.DecodeDir(dir))
		go func(cmd, dir string) {
			defer wg.Done()
			if err := LogCpuRun(cmd, dir); err != nil {
				Err().Println(err)
			}
		}(cmd_str, s.DecodeDir(dir))
	}
	wg.Wait()
	return nil
}

func (s Sat) Score(set string) ([][]string, error) {
	return AutoScore(s.Identify(), DecodeDirs(set, s))
}

func (s Sat) Identify() string {
	return "SAT"
}

type SatTask struct {
	Sat
	TaskConf *TaskConf
}

func NewSatTask() *SatTask {
	return &SatTask{*NewSat(), NewTaskConf()}
}

func (t SatTask) Identify() string {
	return t.Sat.Identify()
}

func (t SatTask) Run() error {
	return Run(t.TaskConf, t.Sat)
}

func SatTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewSatTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("Sat Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
