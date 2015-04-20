//
package main

import (
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"github.com/rosrad/util"
	"strings"
	"sync"
)

func GripSvr() []string {
	grip := []string{}
	exclude := []int{8}
	const MaxNode = 13
	for i := 1; i < MaxNode+1; i++ {
		if util.Contains(exclude, i) {
			continue
		}
		grip = append(grip, fmt.Sprintf("node%02d", i))
	}
	return grip
}

type ProcInfo struct {
	Host     string
	ProcName string
	Id       int
}

func SyncPs(svr string) []ProcInfo {
	procs := []ProcInfo{}
	cmd_str := kaldi.JoinArgs("ssh",
		svr,
		"\"ps -u ren -o comm\"")
	out, err := kaldi.BashOutput(cmd_str)
	if err != nil {
		// log.Println("Err:", err)
	} else {
		str := strings.Trim(string(out), "\n ")
		ps := strings.Split(str, "\n")
		for _, v := range ps {
			procs = append(procs, ProcInfo{svr, v, 1})
		}
	}
	return procs
}
func Exclude(file string) string {
	lines, _ := util.ReadLines(file)
	return strings.Join(lines, " ")
}

func FilterPs(grid_ps []ProcInfo, exclude string) []ProcInfo {
	flt_ps := []ProcInfo{}
	for _, ps := range grid_ps {
		if !strings.Contains(exclude, ps.ProcName) {
			flt_ps = append(flt_ps, ps)
		}
	}
	return flt_ps
}

func PrintPs(ps []ProcInfo) {
	for _, p := range ps {
		fmt.Println(p.Host, p.ProcName)
	}
}

func ParallelSync(exclude_str string) []ProcInfo {
	ps := []ProcInfo{}
	rec := make(chan ProcInfo)
	var wg sync.WaitGroup
	for _, svr := range GripSvr() {
		wg.Add(1)
		go func(svr string) {
			defer wg.Done()
			i := 0
			for _, p := range FilterPs(SyncPs(svr), exclude_str) {
				rec <- p
				i++
			}
			// fmt.Println(i)
		}(svr)
	}

	done := make(chan struct{})
	go func() {
		idx := 0
		for {
			select {
			case p, ok := <-rec:
				if ok {
					ps = append(ps, p)
					idx = idx + 1
					// fmt.Printf("%dth value\n", idx)
				} else {
					// fmt.Println("channel closed with idx =", idx)
					return
				}
			case <-done:
				return
			default:
			}
		}

	}()
	wg.Wait()
	done <- struct{}{}
	close(done)
	close(rec)
	return ps
}

func main() {
	var exculde string
	flag.StringVar(&exculde, "exclude", "", "exclude file for ps grep")
	flag.Parse()
	exclude_str := Exclude(exculde)
	PrintPs(ParallelSync(exclude_str))
}
