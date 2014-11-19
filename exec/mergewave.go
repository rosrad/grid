package main

import (
	"flag"
	"fmt"
	"monc"
	"os"
	"strings"
	"path/filepath"
)

func main() {
	var in_file, outdir string
	var mode int
	flag.StringVar(&in_file, "i", "", "input file ")
	flag.StringVar(&outdir, "o", "", "out dir ")
	flag.IntVar(&mode, "m", 0, "the mode of the program")
	flag.Parse()

	fmt.Println("in_file : ", in_file)
	fmt.Println("outdir : ", outdir)
	fmt.Println("mode :", mode)
	switch mode {
	case 0:
		monc.MergeRaw(in_file, outdir)
	case 1:
		list := strings.Split(in_file, ".")
		list[1] = "raw"
		mer_file := strings.Join(list, ".")
		fmt.Println("merged file : ", mer_file)
		monc.CutRaw(in_file, mer_file, outdir)
	case 2:
		filepath.Walk(in_file, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && filepath.Ext(path) == ".list" {
				dir, _ := filepath.Split(path)
				for i := 0; i < 13; i++ {
					mer_file := fmt.Sprintf("%s/wavech%d.raw", dir, i)
					fmt.Println("Merge File", mer_file)
					splite_dir := fmt.Sprintf("%s/wavech%d/", outdir,i)
					fmt.Println("OutDir", splite_dir)
					os.MkdirAll(splite_dir,os.ModePerm)
					monc.CutRaw(path,mer_file, splite_dir)
				}
			}
			return nil
		})

	}
}












