package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"os"
	"strconv"
	"strings"
)

func TasksFrom(identifier string) []kaldi.TaskRuner {
	file := kaldi.TaskFile(identifier)
	kaldi.Trace().Println("Selected Task:", file)
	fts, err := os.Open(file)
	if err != nil {
		kaldi.Err().Println("Can not open task:", file)
	}
	defer fts.Close()
	task_func := kaldi.TaskFromFunc(identifier)
	if task_func == nil {
		kaldi.Err().Println("No effective TaskFromFunc for ", identifier)
		return []kaldi.TaskRuner{}
	}
	return task_func(fts)
}

func SingleTask(tag string, idx int) {
	tasks := TasksFrom(tag)
	length := len(tasks)
	if idx >= length {
		kaldi.Err().Println("Out of task range :", length)
		return
	}
	if idx >= 0 {
		kaldi.Trace().Println("Task Idx:", idx)
		tasks[idx].Run()
		return
	}
	kaldi.Trace().Printf("%d Tasks in file [%s] will be run!\n", len(tasks), tag)
	for i, t := range tasks {
		kaldi.Trace().Println("Task Idx:", i)
		t.Run()
	}
}

func RunMultiTask(list string) {
	list_file := kaldi.TaskFile(list)
	fl, err := os.Open(list_file)
	if err != nil {
		kaldi.Err().Println("Can not open task list:", list_file)
	}
	defer fl.Close()

	sc := bufio.NewScanner(fl)
	kaldi.Trace().Println("Task List File:", list_file)
	for sc.Scan() {
		str := sc.Text()
		kaldi.Trace().Println("Task String:", str)
		items := strings.Split(str, "#")
		RunTask(items[0], items[1])
	}
}
func parseRange(str string) ([]int, error) {
	kaldi.Trace().Println("Scope string", str)
	res := []int{}
	subpart := strings.Split(str, ",")
	for _, part := range subpart {
		items := strings.Split(part, "~")
		boundary := make([]int, len(items), len(items))
		for idx, _ := range items {
			if len(items[idx]) == 0 {
				boundary[idx] = -1
				continue
			}
			tmp, err := strconv.ParseInt(items[idx], 10, 32)
			if err != nil {
				return []int{}, fmt.Errorf("Syntax Err: parse item to int err: ", err)
			}
			boundary[idx] = int(tmp)
		}

		if len(boundary) > 2 {
			return []int{}, fmt.Errorf("Syntax Err: range item pair disatteched ")
		}
		if len(boundary) == 1 {
			boundary = append(boundary, boundary[0]+1)
		}
		for i := boundary[0]; i < boundary[1]; i++ {
			res = append(res, i)
		}
	}
	return res, nil
}

func RunTask(tag, scope string) error {
	idx_list, err := parseRange(scope)
	if err != nil {
		kaldi.Err().Println("Parse range:", err)
		return fmt.Errorf("Parse range:", err)
	}
	kaldi.Trace().Println("Range:", idx_list)
	for _, value := range idx_list {
		SingleTask(tag, value)
	}
	return nil
}

func main() {
	var list, scope string
	flag.StringVar(&list, "list", "", "task list")
	flag.StringVar(&scope, "scope", "-1", "scope: number ,1-3,4,5 ")
	flag.PrintDefaults()
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("No enough args!")
		return
	}
	tag := flag.Arg(0)
	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("tag", tag)
	kaldi.Trace().Println("task-run")
	kaldi.Trace().Println("Tasks: DNN, GMM, MK-BNF")

	if list == "" {
		RunTask(tag, scope)
	} else {
		RunMultiTask(list)
	}

}
