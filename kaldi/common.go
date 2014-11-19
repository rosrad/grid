//
package kaldi

import (
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

func BashRun(cmd string) error {
	return sh.Command("bash", "-c", cmd).Run()
}

func BashOutput(cmd string) (out []byte, err error) {
	return sh.Command("bash", "-c", cmd).Output()
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
	blacklist := []string{"dt", "decode", "for", "reverb", "simdata"}
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
