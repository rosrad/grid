package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rosrad/kaldi"
	"os"
	"strconv"
	"strings"
	"sync"
)

func TasksFrom(identifier string) []kaldi.TaskRuner {
	file := kaldi.TaskFile(identifier)
	kaldi.Trace().Println("Selected Task:", file)
	fts, err := os.Open(file)
	if err != nil {
		kaldi.Err().Println("Can not open task:", file)
	}
	defer fts.Close()
	taskId := strings.Split(identifier, "_")[0]
	task_func := kaldi.TaskFromFunc(taskId)
	if task_func == nil {
		kaldi.Err().Println("No effective TaskFromFunc for ", taskId)
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
	var wg sync.WaitGroup
	for i, t := range tasks {
		wg.Add(1)
		go func(idx int, t kaldi.TaskRuner) {
			defer wg.Done()
			kaldi.Trace().Println("Task Idx:", idx)
			t.Run()
		}(i, t)
	}
	wg.Wait()
}

func RunMultiTask(list string) {
	list_file := kaldi.TaskList(list)
	fl, err := os.Open(list_file)
	if err != nil {
		kaldi.Err().Println("Can not open task list:", list_file)
	}
	defer fl.Close()

	sc := bufio.NewScanner(fl)
	kaldi.Trace().Println("Task List File:", list_file)
	var wg sync.WaitGroup
	for sc.Scan() {
		str := sc.Text()
		kaldi.Trace().Println("Task String:", str)
		if strings.HasPrefix(str, "//") || len(str) == 0 {
			continue
		}
		wg.Add(1)
		go func(str string) {
			defer wg.Done()
			for _, field := range strings.Split(str, ";") {
				items := strings.Split(field, "#")
				RunTask(strings.Trim(items[0], " \n"), strings.Trim(items[1], " \n"))
			}
		}(str)
	}
	wg.Wait()
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
			boundary = append(boundary, boundary[0])
		}
		for i := boundary[0]; i <= boundary[1]; i++ {
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
	var wg sync.WaitGroup
	for _, value := range idx_list {
		wg.Add(1)
		go func(tag string, idx int) {
			defer wg.Done()
			SingleTask(tag, idx)
		}(tag, value)
	}
	wg.Wait()
	return nil
}

func main() {
	var list, scope, root, lm string
	var manual bool
	flag.StringVar(&root, "root", "", "root path")
	flag.StringVar(&lm, "lm", "", "language model")
	flag.StringVar(&list, "list", "", "task list")
	flag.StringVar(&scope, "scope", "-1", "scope: number ,1-3,4,5 ")
	flag.BoolVar(&manual, "manual", false, "manual to control servers (default=false) ")
	flag.Parse()
	kaldi.Init(root, lm)
	defer kaldi.Uninit()
	kaldi.Trace().Println("task-run")
	if !manual {
		kaldi.DevInstance().AutoSync()
		kaldi.DevInstance().SortGpu()
		kaldi.DevInstance().PrintNodes(true)
		kaldi.DevInstance().SortCpu()
		kaldi.DevInstance().PrintNodes(false)
	}
	if list != "" {
		RunMultiTask(list)
		return
	}

	if flag.NArg() < 1 {
		fmt.Println("No enough args!")
		return
	}
	tag := flag.Arg(0)
	kaldi.Trace().Println("tag", tag)
	kaldi.Trace().Println("Tasks: DNN, GMM, MK-BNF,MK-FMLLR")
	RunTask(tag, scope)
}
