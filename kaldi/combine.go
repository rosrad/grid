package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"sync"
)

type Combine struct {
	Model
	Src2 ExpBase
}

func NewCombine() *Combine {
	c := &Combine{*NewModel(), *NewExpBase()}
	c.Dst.Exp = c.Identify()
	return c
}

func (c Combine) Identify() string {
	return "COMBINE"
}

func (c Combine) TargetDir() string {
	return c.Dst.ExpDir()
}

func (c Combine) Src1Decode(set string) string {
	return MkDecode(c.Src.ExpDir(), set)
}

func (c Combine) Src2Decode(set string) string {
	return MkDecode(c.Src2.ExpDir(), set)
}

func (c Combine) Subsets(set string) ([]string, error) {
	return c.Src.Subsets(set)
}

func (c Combine) DecodeDir(set string) string {
	return MkDecode(c.TargetDir(), set)
}

func (c Combine) OptStr() string {
	return ""
}

func (c Combine) Train() error {
	return nil
}

func (c Combine) Decode(set string) error {
	dirs, err := c.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, dir := range dirs {
		cmd_str := JoinArgs(
			"local/score_combine.sh",
			path.Join(c.Src.DataDir(), dir),
			TestLang(),
			c.Src1Decode(dir),
			c.Src2Decode(dir),
			c.DecodeDir(dir))
		Trace().Println(cmd_str)
		wg.Add(1)
		go func(cmd, dir string) {
			defer wg.Done()
			if err := LogCpuRun(cmd, dir); err != nil {
				Err().Println(err)
			}
		}(cmd_str, c.DecodeDir(dir))
	}
	wg.Wait()
	return nil
}

func (c Combine) Score(set string) ([][]string, error) {
	return AutoScore(c.Identify(), DecodeDirs(set, c))
}

type CombineTask struct {
	Combine
	TaskConf *TaskConf
}

func NewCombineTask() *CombineTask {
	return &CombineTask{*NewCombine(), NewTaskConf()}
}

func (c CombineTask) Identify() string {
	return c.Combine.Identify()
}

func (c CombineTask) Run() error {
	c.TaskConf.Btrain = false // ignore train step for system combination
	c.TaskConf.Bgraph = false // ignore train step for system combination
	return Run(c.TaskConf, c.Combine)
}

func CombineTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewCombineTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
