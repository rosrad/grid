//
package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func file_files(dir, ext string) []string {
	reg := path.Join(dir, "*."+ext)
	fmt.Println("reg", reg)
	files, _ := filepath.Glob(reg)
	return files
}

func deal_scp(dir string) {
	file := file_files(dir, "scp")
	if len(file) < 1 {
		fmt.Println("No scp files, maybe something wrong!")
	}
	for _, f := range file {
		name := path.Base(f)
		parts := strings.Split(name, "_")
		dir := path.Join(path.Dir(f), strings.Join(parts[:len(parts)-1], "_"))
		kaldi.InsureDir(dir)
		os.Rename(f, path.Join(dir, parts[len(parts)-1]))
	}
}

func deal_const(src, dst string) {
	file := file_files(dst, "*")
	if len(file) < 1 {
		fmt.Println("No scp files, maybe something wrong!")
	}

	for _, f := range file {

		name := path.Base(f)
		parts := strings.Split(name, ".")
		dir := path.Join(path.Dir(f), strings.Join(parts[:len(parts)-1], "."))
		if kaldi.DirExist(dir) {
			f_src := path.Join(src, path.Base(dir), parts[len(parts)-1])
			os.Rename(f_src, path.Join(dir, parts[len(parts)-1]))
		}
	}

}

func main() {
	fmt.Println("mk-scp")
	var src string
	var nocopy bool
	flag.StringVar(&src, "src", "", "source of common data dir")
	flag.BoolVar(&nocopy, "nocopy", false, "whether copy the const file")
	flag.Parse()
	fmt.Println("Narg", flag.NArg())
	if flag.NArg() < 1 {
		fmt.Println("No enough parameters")
	}
	dir := flag.Arg(0)
	deal_scp(dir)
	if !nocopy {
		deal_const(src, dir)
	}
}
