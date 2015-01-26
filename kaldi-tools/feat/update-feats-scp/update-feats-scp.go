//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func ConstructScp(scp string) error {
	ark := strings.Replace(scp, ".scp", ".ark", 1)
	if path.Ext(ark) != ".ark" {
		return fmt.Errorf("Error", ark, path.Ext(ark))
	}
	if _, err := os.Stat(ark); os.IsNotExist(err) {
		return fmt.Errorf("Error", "ark file no exist", ark)
	}

	bakup := scp + ".bakup"
	os.Rename(scp, bakup)
	f, err := os.Open(bakup)
	if err != nil {
		return err
	}
	defer f.Close()

	fw, werr := os.Create(scp)
	if werr != nil {
		return werr
	}
	defer fw.Close()
	bfw := bufio.NewWriter(fw)
	re := regexp.MustCompile("([\\w]+\\s+)[^\\s:]*(:.*)")
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ark, err := filepath.Abs(ark)
		if err != nil {
			return err
		}
		txt := re.ReplaceAllString(sc.Text(), "${1}"+ark+"${2}")
		bfw.WriteString(txt + "\n")
	}
	bfw.Flush()
	return nil
}

func Clear(param string) error {
	keys := []string{"*cmvn*", "feats.scp*"}
	for _, key := range keys {
		cmd := kaldi.JoinArgs(
			"find", param, "-type f", "-iname", "\""+key+"\"",
			"|xargs rm")
		fmt.Println(cmd)
		err := kaldi.BashRun(cmd)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func MergeScp(dir string) error {
	file := path.Join(dir, "feats.scp")
	os.Rename(file, file+".bakup")
	cmd := kaldi.JoinArgs(
		"cat", path.Join(dir, "*.scp"),
		"|sort",
		"> ", file)
	return kaldi.BashRun(cmd)
}

func DirCheck(p string, info os.FileInfo, err error) error {
	if info.IsDir() {
		files, _ := filepath.Glob(path.Join(p, "*.scp"))
		if len(files) == 0 {
			return nil
		}
		fmt.Println("Dir", p)
		for _, f := range files {
			if err := ConstructScp(f); err != nil {
				fmt.Println(err)
			}
		}
		MergeScp(p)
	}
	return nil
}

func ConstructData(data, param string) {
	cmd := kaldi.JoinArgs(
		"find", param, "-type f", "-iname \"feats.scp\"")
	output, _ := kaldi.BashOutput(cmd)
	files := strings.Split(strings.Trim(string(output), " \n"), "\n")
	if len(files) == 0 {
		return
	}
	for _, f := range files {
		dst := strings.Replace(f, "param/", "data/", 1)
		kaldi.InsureDir(path.Dir(dst))
		cmd := kaldi.JoinArgs("cp", f, dst)
		fmt.Println(cmd)
		kaldi.BashRun(cmd)
	}
}

func main() {
	fmt.Println("update-feats-scp")
	var data string
	flag.StringVar(&data, "data", "", "data dir")
	flag.Parse()
	param := flag.Arg(0)
	fmt.Println("param", param)
	if !kaldi.DirExist(param) {
		fmt.Println("No exist param dir:", param)
		return
	}
	Clear(param)
	filepath.Walk(param, DirCheck)
	if data != "" {
		ConstructData(data, param)
	}
}
