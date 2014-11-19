package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"util"
)

type Item struct {
	word     string
	ext      string
	meanings []string
}

func fieldFilter(r rune) bool {
	return strings.ContainsRune(";；", r)
}

func cuter(line string) (it Item, suc bool) {
	items := strings.FieldsFunc(line, fieldFilter)
	if len(items) < 2 {
		// fmt.Println(line)
		return it, false
	}

	it.word = strings.Replace(items[0], " ", "", -1)
	ext_start := 1
	for ; ext_start < len(items); ext_start++ {
		if strings.ContainsAny(items[ext_start], "[]【】") {
			it.ext += strings.Trim(items[ext_start], "[]【】><+ \t")
		} else {
			break
		}
	}

	meanings := ""
	for i := ext_start; i < len(items); i++ {
		if strings.ContainsAny(items[i], "0123456789") {
			if len(meanings) != 0 {
				meanings += "-"
			}
			meanings += strings.Trim(items[i], "[]【】 \t\r\n-")
		}
	}
	if len(meanings) == 0 {
		meanings = items[ext_start]
	}
	it.meanings = strings.Split(meanings, "-")
	return it, true
}

func main() {
	var file string
	flag.StringVar(&file, "file", "", "file of the words")
	flag.Parse()
	fmt.Println(file)
	file, _ = filepath.Abs(file)
	outdir := filepath.Join(filepath.Dir(file), "hujiang")
	os.RemoveAll(outdir)
	os.MkdirAll(outdir, os.ModeDir)
	out := filepath.Join(outdir, filepath.Base(file)+".new")
	fout, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	defer fout.Close()
	if err != nil {
		fmt.Println("create output file Error", err)
		return
	}
	lines, err := util.ReadLines(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, line := range lines {
		res, suc := cuter(line)
		if suc {
			newline := res.word + res.ext + "-"
			newline += strings.Join(res.meanings, "<br/>")
			fout.WriteString(newline + "\r\n")
		}
	}
}
