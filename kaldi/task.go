package kaldi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type TaskBase struct {
	Name  string
	Tag   string
	Time  string
	Valid bool
}

func NewTaskBase(name, tag string) *TaskBase {
	return &TaskBase{name, tag, time.Now().Format("15:04:05 02/01/2006"), true}
}

func (t *TaskBase) Run() error {
	if !t.Valid {
		return fmt.Errorf("Unavailable Task [%s]", t.Name)
	}
	Trace().Printf("Run task [%s] firstly at %s", t.Name, t.Time)
	// Run task script
	Trace().Printf("Run task [%s] NOW at %s", t.Name, time.Now().Format("15:20:30 28/12/2006"))
	return nil
}

func TaskFile(name string) string {
	return path.Join(TaskDir(), name) + ".json"
}

func TaskDir() string {
	return path.Join(RootPath(), "task")
}

type FuncTaskFrom func(io.Reader) []TaskRuner

// return tasks from io.reader functions
func TaskFromFunc(identifer string) FuncTaskFrom {
	switch strings.ToLower(identifer) {
	case "mk-bnf":
		return BnfTasksFrom
	case "align":
		return AlignTasksFrom
	case "gmm":
		return GmmTasksFrom
	case "dnn":
		return DnnTasksFrom
	}
	return nil
}

func WriteTask(t IdentifyTaskRuner) (string, error) {
	InsureDir(TaskDir())
	task_file := TaskFile(t.Identify())
	Trace().Println("TaskFile:", task_file)
	fo, ferr := os.OpenFile(task_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer fo.Close()
	if ferr != nil {
		return "", ferr
	}
	en := json.NewEncoder(fo)
	// b, _ := json.MarshalIndent(t)
	// fmt.Println("json:", string(b))
	return task_file, en.Encode(t)
}

type TaskRuner interface {
	Run() error
}
type IdentifyTaskRuner interface {
	TaskRuner
	Identifer
}

type MdlTask interface {
	Trainer
	Decoder
	Counter
	Identifer
}

type TaskConf struct {
	Btrain  bool
	Bdecode bool
	Bscore  bool
}

func NewTaskConf() *TaskConf {
	return &TaskConf{true, true, true}
}

func Run(conf *TaskConf, runer MdlTask) error {
	if conf.Btrain {
		if err := runer.Train(); err != nil {
			return err
		}
		if err := MkGraph(runer.TargetDir()); err != nil {
			return err
		}
	}
	if conf.Bdecode {
		if err := DecodeSets(runer, DataSets(true)); err != nil {
			return err
		}
	}
	if conf.Bscore {
		result := ScoreSets(runer, DataSets(true))
		res_str := fmt.Sprintf("\n#%s\n%s\n",
			runer.TargetDir(),
			FormatScore(result))
		Trace().Println(res_str)
		bf := bytes.NewBufferString(res_str)
		if err := ResultTo(bf, runer.Identify()); err != nil {
			return err
		}
	}
	return nil
}
