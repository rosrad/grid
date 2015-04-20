//
package kaldi

import (
	"encoding/json"
	"io"
	"path"
	"sync"
)

type Fmllr struct {
	Src ExpBase
	Feat
}

func NewFmllr() *Fmllr {
	return &Fmllr{*NewExpBase(), *NewFeat()}
}

func (f Fmllr) Dst() *ExpBase {
	dst := f.Src
	dst.Label = JoinParams(JoinDots(f.Identify(), f.Src.Feat, f.Src.Label), f.Src.Name)
	return &dst
}

func (f Fmllr) TransformDir(set string) string {
	return path.Join(RootPath(), "feats", f.Dst().Feat, f.Dst().Label, "transform", set)
}

func (f Fmllr) Identify() string {
	return "FMLLR"
}

func (f Fmllr) Condition() string {
	cond := "mc"
	if !f.MC {
		cond = "cln"
	}
	return cond
}

func (f Fmllr) TrainMatrix() error {
	// using align_fmllr to get the align fmllr matrix
	cmd_str := JoinArgs(
		"steps/align_fmllr.sh",
		"--nj", MaxNum(f.Src.TrainData(f.Condition())),
		f.Src.TrainData(f.Condition()),
		Lang(),
		f.Src.ExpDir(),
		f.TransformDir(TrainName(f.Condition())))
	Trace().Println(cmd_str)
	if err := BashRun(cmd_str); err != nil {
		Err().Println("TrainMatrix :", err)
		return err
	}
	return nil
}

func (f Fmllr) FakeDecodeDir(dir string) string {
	return path.Join(f.Src.ExpDir(), "fmllr-decode#"+path.Base(dir))
}

func (f Fmllr) TestMatrix(set string) error {
	dirs, err := f.Src.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		cmd_str := JoinArgs(
			"steps/decode_fmllr.sh",
			"--nj", MaxNum(path.Join(f.Src.DataDir(), dir)),
			Graph(f.Src.ExpDir()),
			path.Join(f.Src.DataDir(), dir),
			f.FakeDecodeDir(dir))
		// copy files
		// ensure dst dir
		InsureDir(f.TransformDir(dir))
		cp_cmd := JoinArgs(
			"rsync",
			"-az",
			"--delete-excluded",
			f.FakeDecodeDir(dir)+"/",
			f.TransformDir(dir)+"/")

		rm_cmd := JoinArgs("rm -fr",
			f.FakeDecodeDir(dir),
			f.FakeDecodeDir(dir)+".si")

		go func(cmd_str, cp_cmd, rm_cmd string) {
			defer wg.Done()
			if err := BashRun(cmd_str); err != nil {
				Err().Println("TrainMatrix :", err)
			}
			if err := BashRun(cp_cmd); err != nil {
				Err().Println("Copy matrix :", err)
			}
			if err := BashRun(rm_cmd); err != nil {
				Err().Println("remove no need decode-fmllr files :", err)
			}
		}(cmd_str, cp_cmd, rm_cmd)
	}
	wg.Wait()
	return nil
}

func (f Fmllr) MkMatrix() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f.TrainMatrix(); err != nil {
			Err().Println("Making TrainMatrix:", err)
		}
	}()

	for _, set := range DataSets(DEV) {
		wg.Add(1)
		go func(set string) {
			defer wg.Done()
			if err := f.TestMatrix(set); err != nil {
				Err().Println("Making Matrix:", err)
			}
		}(set)
	}
	wg.Wait()
	return nil
}

func (f Fmllr) TransformAll() error {
	var wg sync.WaitGroup
	for _, set := range DataSets(TRAIN_MC_TEST) {
		wg.Add(1)
		go func(set string) {
			defer wg.Done()
			if err := f.Transform(set); err != nil {
				Err().Println("Fmllr Transform", err)
			}
		}(set)
	}
	wg.Wait()
	return nil
}

func (f Fmllr) Transform(set string) error {
	dirs, err := f.Src.Subsets(set)
	if err != nil {
		return err
	}

	Trace().Println("Transform:", set)
	Trace().Println("Subdirs:", dirs)
	var wg sync.WaitGroup

	for _, dir := range dirs {
		wg.Add(1)
		transform_dir := f.TransformDir(dir)
		cmd_str := JoinArgs(
			"steps/nnet/make_fmllr_feats.sh",
			"--nj", MaxNum(path.Join(f.Src.DataDir(), dir)),
			"--transform-dir", transform_dir,
			path.Join(f.Dst().DataDir(), dir), // target
			path.Join(f.Src.DataDir(), dir),   // source
			f.Src.ExpDir(),
			f.Dst().LogDir(),
			path.Join(f.Dst().ParamDir(), dir))
		go func(cmd_str string) {
			defer wg.Done()
			if err := BashRun(cmd_str); err != nil {
				Err().Println("TrainMatrix :", err)
			}
		}(cmd_str)
	}
	wg.Wait()
	return nil
}

type FmllrTask struct {
	Fmllr
	Bestimation bool
}

func NewFmllrTask() *FmllrTask {
	return &FmllrTask{*NewFmllr(), true}
}

func (f FmllrTask) Identify() string {
	return f.Fmllr.Identify()
}

func (f FmllrTask) Run() error {
	if f.Bestimation {
		if err := f.MkMatrix(); err != nil {
			Err().Println("FMLLR Estimation ", err)
		}
	}

	return f.TransformAll()
}

func FmllrTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewFmllrTask()
		err := dec.Decode(t)
		if err != nil {
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
