package kaldi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rosrad/test/testlib"
	"io"
	"os"
	"path"
	"strings"
)

func TaskList(name string) string {
	return path.Join(TaskListDir(), name) + ".json"
}
func TaskFile(name string) string {
	return path.Join(TaskUnitDir(), name) + ".json"
}

func TaskListDir() string {
	return path.Join(TaskDir(), "list")
}

func TaskUnitDir() string {
	return path.Join(TaskDir(), "unit")
}
func TaskDir() string {
	return path.Join(RootPath(), "task")
}

type FuncTaskFrom func(io.Reader) []TaskRuner

// return tasks from io.reader functions
func TaskFromFunc(identifer string) FuncTaskFrom {

	switch strings.ToLower(identifer) {
	case "bnf":
		return BnfTasksFrom
	case "fmllr":
		return FmllrTasksFrom
	case "align":
		return AlignTasksFrom
	case "cmvn":
		return CmvnTasksFrom
	case "mono":
		return MonoTasksFrom
	case "gmm":
		return GmmTasksFrom
	case "lda":
		return LdaTasksFrom
	case "net":
		return NetTasksFrom
	case "dnn":
		return DnnTasksFrom
	case "discdnn":
		return DiscDnnTasksFrom
	case "splicedgmm":
		return SplicedGmmTasksFrom
	case "ubm":
		return UbmTasksFrom
	case "plda":
		return PldaTasksFrom
	case "combine":
		return CombineTasksFrom
	case "paste":
		return PasterTasksFrom
	case "mfcc":
		return MfccTasksFrom
	}

	return nil
}

func WriteTask(t IdentifyTaskRuner) (string, error) {
	InsureDir(TaskUnitDir())
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
	Bgraph  bool
	Bdecode bool
	Bscore  bool
}

func NewTaskConf() *TaskConf {
	return &TaskConf{true, true, true, true}
}

func Run(conf *TaskConf, runer MdlTask) error {
	dev := SystemSet()
	if conf.Btrain && SysConf().Btrain {
		if err := runer.Train(); err != nil {
			return err
		}
	}
	if conf.Bgraph && SysConf().Bgraph {
		if err := MkGraph(runer.TargetDir()); err != nil {
			return err
		}
		fmt.Println("graph finished")
	}

	if conf.Bdecode && SysConf().Bdecode {
		fmt.Println("decoding")
		if err := DecodeSets(runer, DataSets(dev)); err != nil {
			Err().Println("Decode Err:", err)
			return err
		}
	}
	if conf.Bscore && SysConf().Bscore {
		msg := ""
		for _, set := range DataSets(dev) {
			// result := ScoreSets(runer, set)
			result, err := runer.Score(set)
			if err != nil {
				Trace().Println("Score Set", set)
				continue
			}

			res_str := fmt.Sprintf("#%s\n%s\n",
				runer.TargetDir(),
				FormatScore(result))
			Trace().Println(res_str)
			if len(result) > 0 {
				msg += res_str
			}
			bf := bytes.NewBufferString(res_str)
			if err := ResultTo(bf, runer.Identify(), set); err != nil {
				return err
			}
		}
		if msg != "" {
			testlib.SendGmail(runer.TargetDir(), msg)
		}
	}
	return nil
}
