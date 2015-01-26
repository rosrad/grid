package kaldi

import (
	"encoding/json"
	"io"
	"path"
)

type Lda struct {
	Model
}

func NewLda() *Lda {
	l := &Lda{*NewModel()}
	l.Dst.Exp = l.Identify()
	return l
}

func (l Lda) Identify() string {
	return "LDA"
}

func (l Lda) TargetDir() string {
	return l.Dst.ExpDir()
}

func (l Lda) AlignDir() string {
	return l.Src.AlignDir()
}

func (l Lda) Subsets(set string) ([]string, error) {
	return l.Dst.Subsets(set)
}

func (l Lda) DecodeDir(set string) string {
	return MkDecode(l.TargetDir(), set)
}

func (l Lda) OptStr() string {
	return l.Feat.OptStr()
	// JoinArgs(l.ModelConf.OptStr(), l.LdaConf.OptStr())
}

func (l Lda) Train() error {
	cmd_str := JoinArgs(
		"steps/train_lda_mllt.sh",
		l.OptStr(),
		GaussConf(),
		l.Dst.TrainData(l.Condition()),
		Lang(),
		l.AlignDir(),
		l.TargetDir(),
	)
	err := LogCpuRun(cmd_str, l.TargetDir())
	if err != nil {
		return err
	}
	return nil
}

// implement the Decoder interface
func (l Lda) Decode(set string) error {
	dirs, err := l.Dst.Subsets(set)
	if err != nil {
		return err
	}

	l.Transform = append(l.Transform, path.Join(l.TargetDir(), "final.mat"))
	for _, dir := range dirs {
		cmd_str := JoinArgs(
			"steps/decode.sh",
			"--nj ", JobNum("decode"),
			l.FeatOpt(),
			Graph(l.TargetDir()),
			path.Join(l.Dst.DataDir(), dir),
			l.DecodeDir(dir))
		if err := LogCpuRun(cmd_str, l.DecodeDir(dir)); err != nil {
			Err().Println(err)
		}
	}
	return nil
}

func (l Lda) Score(set string) ([][]string, error) {
	return AutoScore(l.Identify(), DecodeDirs(set, l))
}

type LdaTask struct {
	Lda
	*TaskConf
}

func NewLdaTask() *LdaTask {
	return &LdaTask{*NewLda(), NewTaskConf()}
}

func (t LdaTask) Identify() string {
	return t.Lda.Identify()
}

func (t LdaTask) Run() error {
	return Run(t.TaskConf, t.Lda)
}

func LdaTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewLdaTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("GMM Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
