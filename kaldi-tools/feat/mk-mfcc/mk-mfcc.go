package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func MakeMfcc(files []string, dst string) error {
	dir, _ := filepath.Abs(dst)
	kaldi.InsureDir(dir)
	log := filepath.Join(dir, "log")
	kaldi.InsureDir(log)
	kaldi.BashRun("rm " + dir + "/*.*")
	jobs := make(chan string, 100)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int, job <-chan string) {
			defer wg.Done()
			for f := range jobs {
				wav, _ := filepath.Abs(f)
				cmd := kaldi.JoinArgs("cd ~/rev_kaldi; ./steps/make_mfcc.sh",
					"--nj 8", wav, log, dir)
				fmt.Println(cmd)
				kaldi.BashRun(cmd)
			}
		}(i, jobs)
	}
	for _, f := range files {
		jobs <- f
	}
	close(jobs)
	wg.Wait()

	return nil
}

func main() {
	fmt.Println("mk-mfcc")
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
		if !info.IsDir() {
			return nil
		}
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			files, err := filepath.Glob(dir + "/*.wav")
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
			MakeMfcc(files, dst_dir)
		}(path)
		return nil
	})
	wg.Wait()
}
