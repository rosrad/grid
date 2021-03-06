package kaldi

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"
)

type GlobalConf struct {
	Root        string
	LM          string
	DecodeSet   string
	MaxNum      int
	FmllrDecode bool
	TaskConf
	ExcludeNode []int
}

var g_root string

var g_conf = &GlobalConf{"tmp", "bg", "", 32, false, *NewTaskConf(), []int{}}

func SetRoot(root string) {
	g_root = root
}

func LoadGlobalConf() error {
	root_conf := "root.json"

	f, err := os.Open(root_conf)
	if err != nil {
		Err().Println("Global Conf Read :", err)
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(SysConf())
	if err != nil {
		Err().Println("Global Conf Decode :", err)
		return err
	}
	return nil
}

func SysConf() *GlobalConf {
	return g_conf
}

func RootPath() string {
	return SysConf().Root
}

func Lang() string {
	// return path.Join(RootPath(), "data", "lang")
	return path.Join("data", "lang")
}

func TestLang() string {
	// return path.Join(RootPath(), "data", "lang_test_"+SysConf().LM+"_5k")
	return path.Join("data", "lang_test_"+SysConf().LM+"_5k")
}

func Graph(target string) string {
	return path.Join(target, "graph_"+SysConf().LM+"_5k")
}

func SpkNum(dir string) int {
	cmd := JoinArgs("cat",
		path.Join(dir, "spk2utt"),
		"|wc -l")
	spknum, err := BashOutput(cmd)
	if err != nil {
		Err().Println(err)
	}
	num, _ := strconv.Atoi(strings.TrimSpace(string(spknum[:len(spknum)])))
	return num
}

func MaxNum(dir string) string {
	max := SysConf().MaxNum
	if nspk := SpkNum(dir); nspk < max {
		max = nspk
	}
	return strconv.Itoa(max)
}

func DecodeCmd(exp string) string {
	gmm_decoder := "steps/decode.sh"
	if SysConf().FmllrDecode {
		gmm_decoder = "steps/decode_fmllr.sh"
	}
	switch exp {
	case "GMM":
		return gmm_decoder
	case "LDA":
		return gmm_decoder
	case "SAT":
		return gmm_decoder
	case "MONO":
		return gmm_decoder
	case "DNN":
		return "steps/nnet2/decode.sh"
	case "NET":
		return "steps/nnet/decode.sh"
	default:
		return ""
	}

}
