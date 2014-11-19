package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"
)

func mergeSlice(a, b []string) []string {
	c := make([]string, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}

func main() {
	var command, para string
	var interval, offset int

	flag.StringVar(&command, "c", "", "command")
	flag.StringVar(&para, "p", "", "parameters")
	flag.IntVar(&interval, "i", 10, "interval")
	flag.IntVar(&offset, "o", 0, "offset")
	flag.Parse()
	fmt.Println("command and para :", command, para, interval, offset)
	nt := time.Now()
	var out bytes.Buffer
	cmd := exec.Command(command)
	cmd.Args = mergeSlice(cmd.Args, strings.Fields(para))
	cmd.Stderr = &out
	cmd.Stdout = &out
	run_t := time.Date(nt.Year(), nt.Month(), nt.Day(), nt.Hour(), nt.Minute(), (nt.Second()+interval)/interval*interval+offset, 0, time.Local)
	fmt.Println("Decided Runnig time at ", run_t)
	ticker := time.NewTicker(time.Millisecond * 10)
	go func() {
		for now := range ticker.C {
			res := math.Abs(run_t.Sub(now).Seconds())
			if res <= 0.05 {
				fmt.Println("Now:", now)
				err := cmd.Start()
				if err != nil {
					fmt.Println("command error : ", err)
				}
				ticker.Stop()
				fmt.Println("Ticker stopped")
			}

		}
	}()

	if offset< 0  {
		offset = -offset	
	}
	fmt.Println("offset fixed: ", offset)
	time.Sleep(time.Duration(interval+offset) * time.Second) 
	cmd.Wait()
	fmt.Println("==================Output======================")
	fmt.Println(string(out.Bytes()))
}
