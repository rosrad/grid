package main

import (
	"encoding/json"
	"flag"
	"github.com/rosrad/kaldi"
	"github.com/rosrad/util"
	"io"
	"os"
	"strings"
)

type Config struct {
	Cmd  string
	List string
	Args string
}

func NewConfig() *Config {
	return &Config{"", "", ""}
}

func ParseStream(reader io.Reader) (*Config, error) {
	dec := json.NewDecoder(reader)
	t := NewConfig()
	err := dec.Decode(t)
	return t, err
}

func ParseFile(f string) (*Config, error) {
	fs, err := os.Open(f)
	defer fs.Close()
	if err != nil {
		return nil, err
	}
	return ParseStream(fs)
}

func ListFiles(f string) []string {
	items, err := util.ReadLines(f)
	if err != nil {
		kaldi.Err().Println("List file decoding error:", err)
		return []string{}
	}
	list := []string{}
	for _, c := range items {
		str := strings.Trim(c, " \n")
		if len(str) != 0 {
			list = append(list, str)
		}
	}
	return list
}

func ApplyCmd(c Config) {
	list := ListFiles(c.List)
	for _, item := range list {
		cmd := kaldi.JoinArgs(c.Cmd, c.Args, item)
		err := kaldi.CpuBashRun(cmd)
		if err != nil {
			kaldi.Err().Println("Cmd Err:", cmd, "\t", err)
		}
	}
}

func main() {
	var config string
	var manual bool
	flag.StringVar(&config, "config", "", "config file")
	flag.BoolVar(&manual, "manual", false, "manual to control servers (default=false) ")
	flag.Parse()
	kaldi.Init("", "")
	defer kaldi.Uninit()
	kaldi.Trace().Println("task-run")
	if !manual {
		kaldi.DevInstance().AutoSync()
		kaldi.DevInstance().SortGpu()
		kaldi.DevInstance().PrintNodes(true)
		kaldi.DevInstance().SortCpu()
		kaldi.DevInstance().PrintNodes(false)
	}
	cfg, err := ParseFile(config)
	if err != nil {
		kaldi.Err().Println("Parsing Config Err:", err)
		return
	}
	ApplyCmd(*cfg)
}
