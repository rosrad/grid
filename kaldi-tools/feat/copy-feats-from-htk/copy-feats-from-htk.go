package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func MakeArk(f, dst string) error {
	key := strings.Split(path.Base(f), ".")[0]
	cmd := kaldi.JoinArgs("{ echo", key,
		"; copy-feats --htk-in --binary=false", f, "- ; } >>", dst)
	log.Println(cmd)
	return kaldi.BashRun(cmd)
}

func DealOneDir(files []string, dst string) {
	kaldi.InsureDir(dst)
	kaldi.BashRun("rm " + dst + "/*.*")
	max_len := len(files)
	step := 1000
	begin := 0
	idx := 0
	var wg sync.WaitGroup
	for begin < max_len {
		fs := []string{}
		if begin+step < max_len {
			fs = files[begin : begin+step]
			begin = begin + step
		} else {
			fs = files[begin:]
			begin = max_len
		}
		idx++
		log.Printf("Seg %d :%d", begin, len(fs))
		wg.Add(1)
		file := fmt.Sprintf("%s/feature.%d", dst, idx)
		go func(file string, fs []string) {
			defer wg.Done()
			for _, f := range fs {
				MakeArk(f, file+".txt")
			}
			cmd := kaldi.JoinArgs("copy-feats",
				"ark:"+file+".txt",
				"ark,scp:"+file+".ark,"+file+".scp")
			kaldi.BashRun(cmd)
			kaldi.BashRun("rm " + file + ".txt")
		}(file, fs)
	}
	wg.Wait()
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Println("No enough args!")
		return
	}
	src := flag.Arg(0)
	dst := flag.Arg(1)
	fmt.Println(src, dst)
	if !kaldi.DirExist(src) {
		log.Println("no exist source dir :", src)
		return
	}

	var wg sync.WaitGroup
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			wg.Add(1)
			go func(dir string) {
				defer wg.Done()
				files, err := filepath.Glob(dir + "/*.mfc")
				if err != nil {
					log.Println(err)
					return
				}
				if len(files) == 0 {
					return
				}
				rel, rerr := filepath.Rel(src, dir)
				if rerr != nil {
					log.Println(err)
					return
				}
				dst_dir := filepath.Join(dst, rel)
				DealOneDir(files, dst_dir)
			}(path)
		}
		return nil
	})
	wg.Wait()
}
