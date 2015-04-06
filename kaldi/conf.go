package kaldi

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
)

type GlobalConf struct {
	Root      string
	LM        string
	DecodeSet string
	TaskConf
	ExcludeNode []int
}

var g_root string

var g_conf = &GlobalConf{"tmp", "bg", "", *NewTaskConf(), []int{}}

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
	return path.Join(RootPath(), "data", "lang")
}

func TestLang() string {
	return path.Join(RootPath(), "data", "lang_test_"+SysConf().LM+"_5k")
}

func Graph(target string) string {
	return path.Join(target, "graph_"+SysConf().LM+"_5k")
}

func JobNum(job string) string {
	parallel := 1
	switch job {
	case "decode":
		parallel = 4
	case "train":
		parallel = 16
	case "dnn":
		parallel = 16
	}
	return strconv.Itoa(parallel)
}

func GaussConf() string {
	return JoinArgs("2500", "15000")
}
