//
package kaldi

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

func JoinParams(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return strings.Trim(strings.Join(elem[i:], "_"), "_")
		}
	}
	return ""
}

func JoinArgs(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return strings.Trim(strings.Join(elem[i:], " "), " ")
		}
	}
	return ""
}

func JoinDots(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return strings.Trim(strings.Join(elem[i:], "."), ".")
		}
	}
	return ""
}

func FileExist(f string) error {
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("file exist error: %s", f)
}

func InsureDir(dir string) {
	if !DirExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func DirExist(dir string) bool {
	if src, err := os.Stat(dir); os.IsNotExist(err) || !src.IsDir() {
		return false
	}
	return true
}

func LogGpuRun(cmd, dir string) error {
	dev_cmd := cmd
	if DevInstance().Inited() {
		opt := DevInstance().AutoSelectGpu()
		cwd, cerr := os.Getwd()
		if cerr != nil {
			return cerr
		}
		dev_cmd = JoinArgs("ssh", opt.Node,
			"\"cd", cwd+";",
			cmd+"\"")
		Trace().Printf("HOST: %s\nCMD: %s\n", opt.Node, dev_cmd)
	}

	return LogRun(dev_cmd, dir)
}

func LogCpuRun(cmd, dir string) error {
	dev_cmd := cmd
	if DevInstance().Inited() {
		opt := DevInstance().AutoSelectCpu()
		cwd, cerr := os.Getwd()
		if cerr != nil {
			return cerr
		}
		dev_cmd = JoinArgs("ssh", opt.Node,
			"\"cd", cwd+";",
			cmd+"\"")
		Trace().Printf("HOST: %s\nCMD: %s\n", opt.Node, dev_cmd)
	}
	return LogRun(dev_cmd, dir)
}

func LogRun(cmd, dir string) error {
	InsureDir(dir)
	file := path.Join(dir, "cmd")
	f, err := os.Create(file)
	if err != nil {
		Err().Println(err)
		return err
	}
	defer f.Close()
	content := cmd + "\n"
	if _, err := f.WriteString(content); err != nil {
		Err().Println(err)
		return err
	}
	tag := path.Base(dir)
	s := sh.Command("bash", "-c", cmd)
	s.Stderr = NewLogWriter(tag)
	s.Stdout = NewLogWriter(tag)
	return s.Run()
}

func CpuBashRun(cmd string) error {
	opt := DevInstance().AutoSelectCpu()
	cwd, cerr := os.Getwd()
	if cerr != nil {
		return cerr
	}
	dev_cmd := JoinArgs("ssh", opt.Node,
		"\"cd", cwd+";",
		cmd+"\"")
	return BashRun(dev_cmd)
}

func BashRun(cmd string) error {
	s := sh.Command("bash", "-c", cmd)
	// s.Stdout = LogWriter()
	return s.Run()
}

func CpuBashOutput(cmd string) (out []byte, err error) {
	opt := DevInstance().AutoSelectCpu()
	cwd, cerr := os.Getwd()
	if cerr != nil {
		return []byte{}, cerr
	}
	dev_cmd := JoinArgs("ssh", opt.Node,
		"\"cd", cwd+";",
		cmd+"\"")
	return BashOutput(dev_cmd)
}

func BashOutput(cmd string) (out []byte, err error) {
	s := sh.Command("bash", "-c", cmd)
	s.Stderr = nil
	return s.Output()
}

func Now() string {
	return time.Now().Format("2006-Jan-02")
}

func Data2Param(data string) string {
	re := regexp.MustCompile("(.*/feats/[^/]+)/(data)/(.*)")
	subpart := re.FindStringSubmatch(data)
	return path.Join(subpart[1], "param", subpart[3])
}

func Contains(item string, set []string) bool {
	for _, value := range set {
		if value == item {
			return true
		}
	}
	return false
}

func Exclude(sets, blacklist []string) []string {
	modified := []string{}
	for _, item := range sets {
		if !Contains(item, blacklist) {
			modified = append(modified, item)
		}
	}
	return modified
}

func ExcludeDefault(set []string) []string {
	blacklist := []string{"dt", "decode", "for", "reverb", "simdata", "1ch", "a", "bg"}
	return Exclude(set, blacklist)
}

func Unique(set string, ignore_case bool) []string {
	var sub sort.StringSlice

	if ignore_case {
		set = strings.ToLower(set)
	}
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	sub = strings.FieldsFunc(set, f)
	sub.Sort()
	prev := ""
	unique := []string{}
	for _, value := range sub {
		if prev != value {
			unique = append(unique, value)
		}
		prev = value
	}
	return unique
}

func Map2Slice(src map[string]string) []string {
	dst := []string{}
	for _, value := range src {
		dst = append(dst, value)
	}
	sort.Sort(sort.StringSlice(dst))
	return dst
}
