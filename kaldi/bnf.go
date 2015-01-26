package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
	"sync"
)

type BnfConf struct {
	BottleDim int
}

func NewBnfConf() *BnfConf {
	return &BnfConf{BottleDim: 42}
}

func (bc BnfConf) OptStr() string {
	var_opt := ""
	var_opt = JoinArgs(var_opt, "--bottleneck-dim", strconv.Itoa(bc.BottleDim))
	return var_opt
}

type BnfDnn struct {
	Dnn // inherit from normal Dnn
	BnfConf
}

func NewBnfDnn() *BnfDnn {
	b := &BnfDnn{*NewDnn(), *NewBnfConf()}
	b.Dst.Feat = "bnf"
	return b
}

func (b BnfDnn) OptStr() string {
	return JoinArgs(b.BnfConf.OptStr(),
		b.DnnConf.OptStr(),
		b.Feat.OptStr())
}

func (b *BnfDnn) Train() error {
	cmd_str := JoinArgs(
		"steps/nnet2/train_tanh_bottleneck.sh",
		" --stage -100",
		"--mix-up 5000",
		"--max-change 40",
		"--initial-learning-rate 0.005",
		"--final-learning-rate 0.0005",
		b.OptStr(),
		b.Src.TrainData(b.Condition()), // here mfcc feature is used to train bnfdnn
		Lang(),
		b.AlignDir(),
		b.TargetDir())

	Trace().Println(cmd_str)
	err := LogGpuRun(cmd_str, b.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

func (b *BnfDnn) TargetDir() string {
	return b.Dst.DeriveExp()
}

func (b *BnfDnn) CleanStorage() bool {
	return b.Dst.Bakup() != 3
}

func (b *BnfDnn) InsureStorage() {
	b.CleanStorage()
	dirs := []string{b.TargetDir(), b.Dst.DataDir(), b.Dst.ParamDir()}
	for _, dir := range dirs {
		InsureDir(dir)
	}
}

func (b *BnfDnn) Dump(set string) error {
	dirs, err := b.Src.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		dst_data := path.Join(b.Dst.DataDir(), dir)
		cmd_str := JoinArgs(
			"steps/nnet2/dump_bottleneck_features.sh",
			"--nj", JobNum("decode"),
			b.FeatOpt(),
			path.Join(b.Src.DataDir(), dir), // source dir
			dst_data,                        // bnf destination dir
			b.TargetDir(),                   // bnf model
			path.Join(b.Dst.ParamDir(), dir),
			b.Dst.DeriveDump())
		go func(cmd, dir string) {
			defer wg.Done()
			if err := LogCpuRun(cmd, dir); err != nil {
				Err().Println("Dump#", err)
			}
		}(cmd_str, dst_data)
	}
	wg.Wait()
	return nil
}

func (b *BnfDnn) DumpSets(sets []string) {
	b.InsureStorage()
	c := NewCmvn()
	c.ExpBase = b.Dst
	var wg sync.WaitGroup
	for _, set := range sets {
		wg.Add(1)
		go func(set string) {
			defer wg.Done()
			b.Dump(set)
			c.Compute(set)
		}(set)
	}
	wg.Wait()
}

type BnfTask struct {
	BnfDnn
	*TaskConf
}

func NewBnfTask() *BnfTask {
	return &BnfTask{*NewBnfDnn(), NewTaskConf()}
}

func (t BnfTask) Identify() string {
	return "BNF"
}

func (t BnfTask) Run() error {
	if t.Btrain {
		if err := t.Train(); err != nil {
			return err
		}
	}
	if t.Bdecode {
		sets := DataSets(DEV_TRAIN)
		Trace().Println("Dataset:")
		Trace().Println(sets)
		t.DumpSets(sets)
	}
	return nil
}

func BnfTasksFrom(reader io.Reader) []TaskRuner {

	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewBnfTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("BNF task decode:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
