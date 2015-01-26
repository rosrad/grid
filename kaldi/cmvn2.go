//
package kaldi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
)

type Cmvn struct {
	ExpBase
}

func NewCmvn() *Cmvn {
	return &Cmvn{*NewExpBase()}
}

func fileExist(f string) error {
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("file exist error: %s", f)
}

func InsureScps(dir string) error {
	utt2utt := path.Join(dir, "utt2utt")
	utt2spk := path.Join(dir, "utt2spk")
	if err := fileExist(utt2spk); err != nil {
		return err
	}
	if err := fileExist(utt2utt); err != nil {
		cmd := "cat " + utt2spk + "| awk '{print $1,$1}' > " + utt2utt
		Trace().Println("Make utt2utt:", utt2utt)
		return BashRun(cmd)
	}
	return nil
}

func (c Cmvn) Compute(set string) error {
	dirs, err := c.Subsets(set)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			if err := InsureScps(path.Join(c.DataDir(), dir)); err != nil {
				Err().Println(err)
				return
			}
			scps := []string{"utt2utt", "spk2utt"}
			for _, scp := range scps {
				param_f := path.Join(c.ParamDir(), dir, JoinParams("cmvn", scp))
				data_f := path.Join(c.DataDir(), dir, JoinParams("cmvn", scp))
				scp_f := path.Join(c.DataDir(), dir, scp)
				comp_cmd := JoinArgs(
					"compute-cmvn-stats",
					"--spk2utt=ark:"+scp_f,
					"scp:"+path.Join(c.DataDir(), dir, "feats.scp"),
					"ark,scp:"+param_f+".ark"+","+param_f+".scp")
				// copy scp from param to data
				cp_cmd := JoinArgs("cp", param_f+".scp", data_f+".scp")

				if err := BashRun(comp_cmd); err != nil {
					Err().Println(err)
					continue
				}
				if err := lenDiff(param_f+".scp", scp_f); err != nil {
					Err().Println(err)
					continue
				}
				BashRun(cp_cmd)
			}
		}(dir)
	}
	wg.Wait()
	return nil
}

func (c Cmvn) ComputeAll() {
	var wg sync.WaitGroup
	for _, set := range DataSets(ALL) {
		wg.Add(1)
		go func(set string) {
			defer wg.Done()
			c.Compute(set)
		}(set)
	}
	wg.Wait()
}

func lenDiff(f1, f2 string) error {
	o1, _ := BashOutput("cat " + f1 + "| wc -l")
	o2, _ := BashOutput("cat " + f2 + "| wc -l")
	if !bytes.Equal(bytes.TrimSpace(o1), bytes.TrimSpace(o2)) {
		return fmt.Errorf("Count no matched %s(%s) != %s(%s)", f1, o1, f2, o2)
	}
	return nil
}

type CmvnTask struct {
	Cmvn
}

func NewCmvnTask() *CmvnTask {
	return &CmvnTask{*NewCmvn()}
}

func (t CmvnTask) Identify() string {
	return "CMVN"
}

func (c CmvnTask) Run() error {
	c.ComputeAll()
	return nil
}

func CmvnTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewCmvnTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("Cmvn Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
