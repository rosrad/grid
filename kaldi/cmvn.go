package kaldi

import (
	"bytes"
	"fmt"
	"os"
	"path"
)

type CmvnOption struct {
	ExpBase
	mode string // for group type
}

func NewCmvnOption(feat, mode string) *CmvnOption {
	c := &CmvnOption{*NewExpBase(), mode}
	c.Feat = feat
	return c
}

func (opt *CmvnOption) HeapName() string {
	dic := map[string]string{
		"utt":  "utt2utt",
		"spk":  "spk2utt",
		"fake": "fake"}
	return dic[opt.mode]
}

func (opt *CmvnOption) CheckFiles(data string) error {
	names := []string{"feats.scp", "spk2utt", "utt2spk"}
	for _, name := range names {
		f := path.Join(data, name)
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return fmt.Errorf("No Exist File: %s", f)
		}
	}
	return nil
}

func (opt *CmvnOption) InsureUtt(data string) {
	utt := path.Join(data, "utt2utt")
	utt2spk := path.Join(data, "utt2spk")
	if _, err := os.Stat(utt); !os.IsNotExist(err) {
		// file exist
		return
	}
	// make the "utt2utt" file
	cmd := "cat " + utt2spk + "| awk '{print $1,$1}' > " + utt
	// fmt.Println("Cmd:", cmd)
	Trace().Println("Make utt2utt:", utt)
	BashRun(cmd)
}

func (opt *CmvnOption) FinalCheck(com1, com2 string) error {
	no, _ := BashOutput("cat " + com1 + "| wc -l")
	nt, _ := BashOutput("cat " + com2 + "| wc -l")
	if !bytes.Equal(bytes.TrimSpace(no), bytes.TrimSpace(nt)) {
		return fmt.Errorf("Count no matched %s != %s", no, nt)
	}
	return nil
}

func (opt *CmvnOption) Paths(subset string) (data, param string) {
	data = path.Join(opt.DataDir(), subset)
	param = path.Join(opt.ParamDir(), subset)
	return
}

func (opt *CmvnOption) CmvnFiles(subset string) (data_file, param_file, heap_file string) {
	data, param := opt.Paths(subset)
	heap := ""
	switch opt.mode {
	// case "fake":
	case "utt":
		heap = "utt2utt"
		opt.InsureUtt(data)
	case "spk":
		heap = "spk2utt"
	default:
		fmt.Errorf("No supported mode :%s", opt.mode)
	}
	name := path.Base(data)
	param_file = path.Join(param, JoinParams("cmvn", heap, name))
	data_file = path.Join(data, JoinParams("cmvn", heap))
	heap_file = path.Join(data, heap)
	return
}

func (opt *CmvnOption) CompCmvn(subset string) error {
	data, _ := opt.Paths(subset)
	if err := opt.CheckFiles(data); err != nil {
		return err
	}
	data_file, param_file, heap_file := opt.CmvnFiles(subset)

	comp_cmd := JoinArgs(
		"compute-cmvn-stats",
		"--spk2utt=ark:"+heap_file,
		"scp:"+path.Join(data, "feats.scp"),
		"ark,scp:"+param_file+".ark"+","+param_file+".scp")
	if err := BashRun(comp_cmd); err != nil {
		return err
	}

	if err := opt.FinalCheck(heap_file, param_file+".scp"); err != nil {
		return err
	}

	// copy the scp file from "param" directory to "data" directory
	cp_cmd := JoinArgs("cp", param_file+".scp", data_file+".scp")
	BashRun(cp_cmd)
	return nil
}
