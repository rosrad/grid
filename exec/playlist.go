package main

import (
	"bufio"
	"flag"
	"fmt"
	"monc"
	"os"
	"os/exec"
	"util"
)

func main() {
	var inlist, outlist, cmd string
	flag.StringVar(&inlist, "i", "", "list of the training data keyword")
	flag.StringVar(&outlist, "o", "", "path list of the training data keyword")
	flag.StringVar(&cmd, "c", "", "path of command")
	flag.Parse()

	if !util.IsExist(outlist) {
		if !monc.GenerateList(inlist, outlist) {
			fmt.Println("playlist generation failed")
		}
	}
	ofr, err := os.Open(outlist)
	if err != nil {
		return
	}
	defer ofr.Close()
	ofrb := bufio.NewScanner(ofr)
	fmt.Println("Start the task")
	for ofrb.Scan() {
		command := exec.Command(cmd, ofrb.Text())
		out,_ := command.CombinedOutput()
		fmt.Println("Run Cmd : ", cmd, ofrb.Text())
		command.Run()
		fmt.Println("Out :", string(out))
	}
	fmt.Println("Finished the task")
}


