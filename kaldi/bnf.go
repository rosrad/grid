package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"strconv"
)

type BnfDnn struct {
	Dnn // inherit from normal Dnn
}

func NewBnfDnn() *BnfDnn {
	b := &BnfDnn{*NewDnn("mfcc")}
	b.Dst.Feat = "bnf"
	return b
}
func (b *BnfDnn) Train() error {
	cmd_str := JoinArgs(
		"steps/nnet2/train_tanh_bottleneck.sh",
		" --stage -100",
		"--num-threads 1 ",
		"--mix-up 5000",
		"--max-change 40",
		"--initial-learning-rate 0.005",
		"--final-learning-rate 0.0005",
		"--bottleneck-dim 42",
		"--hidden-layer-dim 1024",
		b.DnnConf.OptStr(),
		b.ModelConf.OptStr(),
		b.Src.TrainData("mc"), // here mfcc feature is used to train bnfdnn
		Lang(),
		b.AlignDir(),
		b.TargetDir())

	Trace().Println(cmd_str)
	err := BashRun(cmd_str)
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

	for _, dir := range dirs {
		cmd_str := JoinArgs(
			"steps/nnet2/dump_bottleneck_features.sh",
			"--nj "+strconv.Itoa(JobNum("decode")),
			b.ModelConf.OptStr(),
			path.Join(b.Src.DataDir(), dir), // source dir
			path.Join(b.Dst.DataDir(), dir), // bnf destination dir
			b.TargetDir(),                   // bnf model
			path.Join(b.Dst.ParamDir(), dir),
			b.Dst.DeriveDump())

		Trace().Println(cmd_str)
		err := BashRun(cmd_str)
		if err != nil {
			Err().Println("Dump#", err)
		}
	}
	return nil
}

func (b *BnfDnn) DumpSets(sets []string) {
	b.InsureStorage()
	for _, set := range sets {
		b.Dump(set)
		// compute the cmvn status for BNF
		b.CmvnCompute(set)
	}
}

func (b *BnfDnn) CmvnCompute(set string) {
	// compute  cmvn for BNF
	cmvn_modes := []string{"utt", "spk"}
	for _, mode := range cmvn_modes {
		cmvn := NewCmvnOption("bnf", mode)
		cmvn.ExpBase = b.Dst
		subsets, _ := b.Dst.Subsets(set)
		for _, subset := range subsets {
			cmvn.CompCmvn(subset)
		}
	}
}

func (b *BnfDnn) CmvnSets(sets []string) {
	// compute  cmvn for BNF
	for _, set := range sets {
		b.CmvnCompute(set)
	}
}

type BnfTask struct {
	BnfDnn
	TaskBase *TaskBase
}

func NewBnfTask() *BnfTask {
	return &BnfTask{*NewBnfDnn(), NewTaskBase("MK-BNF", "")}
}

func (t BnfTask) Identify() string {
	return "MK-BNF"
}

func (t BnfTask) Run() error {
	if err := t.Train(); err != nil {
		return err
	}
	sets := DataSets(false)
	Trace().Println("Dataset:")
	Trace().Println(sets)
	t.DumpSets(sets)
	return nil
}

func BnfTasksFrom(reader io.Reader) []TaskRuner {

	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewBnfTask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
