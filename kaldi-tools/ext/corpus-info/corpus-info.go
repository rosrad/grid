//
package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"strconv"
	"strings"
)

func main() {

	var src, ext, reg string
	flag.StringVar(&src, "src", "./", "Source directory of the corpus")
	flag.StringVar(&ext, "ext", "wav", "extension  of the corpus with default (wav)")
	flag.StringVar(&reg, "reg", "", "regular expression for searching files)")
	flag.Parse()
	fmt.Printf("src:%s\next:%s\nreg:%s\n", src, ext, reg)
	if !kaldi.DirExist(src) {
		fmt.Printf("Error: No Exist src directory: %s\n", src)
		return
	}
	pattern := "*" + reg + "*." + ext

	file_cmd := kaldi.JoinArgs(
		"find",
		src+"/",
		"-type f ",
		"-iname "+pattern)
	fmt.Println("Find: ", file_cmd)
	output, _ := kaldi.BashOutput(file_cmd)
	files := strings.Split(string(output[:]), "\n")
	total_len := 0.0
	count := 0
	for _, file := range files {
		cmd := kaldi.JoinArgs(
			"sox",
			strings.Trim(file, "\n "),
			"-n stat",
			"2>&1",
			"|awk '/Length/ {print $3}'")
		output, err := kaldi.BashOutput(cmd)
		if err != nil {
			fmt.Println("Missed file:", file)
			continue
		}
		length, err := strconv.ParseFloat(strings.Trim(string(output[:]), "\n \t"), 64)
		if err != nil {
			fmt.Println("Missed file:", file)
			continue
		}
		count++
		total_len += length
	}

	fmt.Println("Total Count: ", len(files))
	fmt.Println("Total Len(hours): ", total_len/3600)
	fmt.Println("Wave Count: ", count)
	fmt.Println("Missed Count: ", len(files)-count)
}
