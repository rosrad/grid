package kaldi

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
)

type Container struct {
	ExpBase
	Feat
}

func NewContainer() *Container {
	return &Container{*NewExpBase(), *NewFeat()}
}

type Paster struct {
	Dst        ExpBase
	Containers []Container
}

func NewPaster() *Paster {
	return &Paster{*NewExpBase(), []Container{}}
}

func (p Paster) Identify() string {
	return "PASTE"
}

// just for implement the interface of decoder
func (p Paster) Decode(set string) error {
	return p.PasteSet(set)
}

func (p Paster) Paste(rel string) error {
	// first step : copy shared files
	src := path.Join(p.Containers[0].DataDir(), rel)
	data := path.Join(p.Dst.DataDir(), rel)
	param := path.Join(p.Dst.ParamDir(), rel)
	InsureDir(data)
	InsureDir(param)
	sync_cmd := JoinArgs("find ", src+"/*", "-maxdepth 0", "-type f",
		"| grep -v cmvn",
		"| xargs -I {} cp {} ", data)
	BashRun(sync_cmd)
	// second step : merge feats.scp

	output := fmt.Sprintf("ark,scp:%s/feature.ark,%s/feature.scp", param, param)
	input := ""
	for _, c := range p.Containers {
		cur_dir := path.Join(c.DataDir(), rel)
		feat := path.Join(cur_dir, "feats.scp")
		if FileExist(feat) != nil {
			err := fmt.Errorf("feats.scp don't exist in Container :%s", cur_dir)
			Err().Println(err)
			return err
		}
		feat_scp := strings.Replace(c.FeatStr(), JobStr(), path.Join(c.DataDir(), rel), -1)
		input = JoinArgs(input, feat_scp)
	}

	paste_cmd := JoinArgs("paste-feats",
		input,
		output)
	LogRun(paste_cmd, p.Dst.DataDir())
	cp_feat := JoinArgs("cp ", path.Join(param, "feature.scp"), path.Join(data, "feats.scp"))
	BashRun(cp_feat)
	return nil
}

func (p Paster) PasteSet(set string) error {
	if len(p.Containers) < 2 {
		err := fmt.Errorf("Container is not enough, just ", len(p.Containers))
		Err().Println(err)
		return err
	}

	dirs, err := p.Containers[0].Subsets(set)
	if err != nil {
		Err().Println(err)
		return err
	}

	for _, dir := range dirs {
		workdir := path.Join(p.Dst.DataDir(), dir)
		InsureDir(workdir)
		p.Paste(dir)
	}
	return nil
}

type PasterTask struct {
	Paster
}

func NewPasterTask() *PasterTask {
	return &PasterTask{*NewPaster()}
}

func (p PasterTask) Identify() string {
	return p.Paster.Identify()
}

func (p PasterTask) Run() error {
	decode_set := TRAIN_MC_TEST
	if len(SysConf().DecodeSet) != 0 {
		decode_set = Str2DataType(SysConf().DecodeSet)
	}
	return DecodeSets(p.Paster, DataSets(decode_set))
}

func PasterTasksFrom(reader io.Reader) []TaskRuner {
	dec := json.NewDecoder(reader)
	tasks := []TaskRuner{}
	for {
		t := NewPasterTask()
		err := dec.Decode(t)
		if err != nil {
			Err().Println("Paster Decode Error:", err)
			break
		}
		tasks = append(tasks, *t)
	}
	return tasks
}
