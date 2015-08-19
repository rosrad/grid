package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/rosrad/kaldi"
	"io"
	"os"
	"strings"
	"sync"
)

type Config struct {
	Cmd  string
	List string
	Args string
}

func NewConfig() *Config {
	return &Config{"", "", ""}
}

func ParseCFGStream(reader io.Reader) []Config {
	cfg := []Config{}
	for {
		dec := json.NewDecoder(reader)
		t := NewConfig()
		err := dec.Decode(t)
		kaldi.Trace().Println("config", t)
		if err != nil {
			kaldi.Err().Println("Parsing Config Error:", err)
			break
		}
		cfg = append(cfg, *t)
	}
	kaldi.Trace().Println("configs", cfg)
	return cfg
}

func ParseCFGFile(f string) []Config {
	fs, err := os.Open(f)
	defer fs.Close()

	if err != nil {
		kaldi.Err().Println("Open config file :", err)
		return []Config{}
	}
	return ParseCFGStream(fs)
}

func ApplyCmd(c Config) {
	file, err := os.Open(c.List)
	kaldi.Trace().Println("list :", c.List)
	if err != nil {
		kaldi.Err().Println("Open list file :", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := strings.Trim(scanner.Text(), " \n")
		if len(item) != 0 {
			cmd := kaldi.JoinArgs(c.Cmd, c.Args, item)
			err := kaldi.CpuBashRun(cmd)
			if err != nil {
				kaldi.Err().Println("Cmd Err:", cmd, "\t", err)
			}
		}
	}
}

func ApplyCmds(cs []Config, max int) {
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func(cfg Config) {
			defer wg.Done()
			ApplyCmd(cfg)
		}(c)
	}
	wg.Wait()
}

func main() {
	var manual bool
	var num int
	flag.BoolVar(&manual, "manual", false, "manual to control servers (default=false) ")
	flag.IntVar(&num, "n", 4, "number of parallel processing")
	flag.Parse()
	kaldi.Init("", "")
	defer kaldi.Uninit()
	kaldi.Trace().Println("cmd-run")

	if flag.NArg() == 0 {
		kaldi.Trace().Println("No enough args!")
		return
	}

	config := flag.Arg(0)
	kaldi.Trace().Println("The first Arg:", config)

	cfgs := ParseCFGFile(config)
	kaldi.Trace().Println(cfgs)
	ApplyCmds(cfgs, num)
}
