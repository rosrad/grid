//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func DirExist(dir string) bool {
	if src, err := os.Stat(dir); os.IsNotExist(err) || !src.IsDir() {
		return false
	}
	return true
}

func TestDirs(src, dst string) bool {
	dirs := []string{path.Join(src, "data"), path.Join(src, "param"), dst}
	for _, dir := range dirs {
		if !DirExist(dir) {
			fmt.Println("Dir Not Exist:", dir)
			return false
		}
	}
	return true
}

func BashRun(cmd string) error {
	return sh.Command("bash", "-c", cmd).Run()
}

func BashOutput(cmd string) (out []byte, err error) {
	return sh.Command("bash", "-c", cmd).Output()
}

func CopyDirs(src, dst string) error {
	source := []string{path.Join(src, "data"), path.Join(src, "param")}
	for _, dir := range source {
		cmd_str := strings.Join([]string{"rsync -axz --delete-excluded ",
			dir, dst + "/"}, " ")
		fmt.Println("cmd :", cmd_str)
		if err := BashRun(cmd_str); err != nil {
			return err
		}
	}
	return nil
}

func FindScp(dir string) []string {
	cmd_str := strings.Join([]string{
		"find",
		dir,
		"-type f",
		"-iname \"*.scp\"",
		"|grep -v wav.scp"}, " ") // exclude the wave.scp
	fmt.Println(cmd_str)
	outstr, _ := BashOutput(cmd_str)
	files := strings.Split(strings.Trim(string(outstr), "\n"), "\n")
	return files
}

func ReplaceContent(file, dst string) error {
	bakup := file + ".bakup"
	os.Rename(file, bakup)

	f, err := os.Open(bakup)
	if err != nil {
		return err
	}
	defer f.Close()

	fw, werr := os.Create(file)
	if werr != nil {
		return werr
	}
	defer fw.Close()

	bfw := bufio.NewWriter(fw)
	re := regexp.MustCompile("([\\w]+\\s+)(\\S+)/param[^/]*/(.*)")
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		dir := dst
		if !filepath.IsAbs(dst) {
			pwd, _ := BashOutput("pwd")
			dir = path.Join(strings.Trim(string(pwd), "\n"), dst)
		}

		prune := re.ReplaceAllString(sc.Text(), "${1}"+dir+"/param/${3}")
		bfw.WriteString(prune + "\n")
	}
	bfw.Flush()
	return nil
}

func PruneParam(src, dst string) {
	subdirs, err := ioutil.ReadDir(src)
	if err != nil || len(subdirs) == 0 {
		fmt.Println("Cannot read any subdirs in Src: ", src)
		return
	}
	for _, dir := range subdirs {
		files := FindScp(path.Join(src, dir.Name()))
		if len(files) == 0 {
			fmt.Println("Empty destination list!", dst)
			return
		}
		fmt.Printf("Files in [%s]:%d\n", src, len(files))
		for _, f := range files {
			file, _ := filepath.Abs(f)
			ReplaceContent(file, dst)
		}
	}
}

func main() {
	fmt.Println("prune-feats")
	var src, dst string
	var bnocopy bool
	flag.StringVar(&src, "src", "./", "the source directory of feature ")
	// flag.StringVar(&dst, "dst", "", "the destination directory of feature ")
	flag.BoolVar(&bnocopy, "nocopy", false, "whether copy source to destination")
	flag.Parse()
	dst = flag.Arg(0)
	fmt.Println("Src:", src)
	fmt.Println("Dst:", dst)
	fmt.Println("NoCopy:", bnocopy)
	fmt.Println("###########################")

	if !bnocopy {
		if !TestDirs(src, dst) {
			return
		}

		fmt.Println("Copy", src, "To", dst)
		if err := CopyDirs(src, dst); err != nil {
			fmt.Println("Copy Error:", err)
		}
	}

	dirs := []string{path.Join(dst, "data"), path.Join(dst, "param")}
	for _, dir := range dirs {
		PruneParam(dir, dst)
	}

}
