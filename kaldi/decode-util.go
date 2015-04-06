//
package kaldi

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
)

type Decoder interface {
	Decode(string) error
}

type Counter interface {
	Score(string) ([][]string, error)
}

type Trainer interface {
	Train() error
	TargetDir() string
}

type Identifer interface {
	Identify() string
}

type DecodeDirMaker interface {
	Subsets(string) ([]string, error)
	DecodeDir(string) string
}

func MkDecode(target, set string) string {
	return path.Join(target, JoinParams("decode", SysConf().LM, "#"+path.Base(set)))
}

func DecodeDirs(set string, mk DecodeDirMaker) []string {
	items, err := mk.Subsets(set)
	dirs := []string{}
	if err != nil {
		Err().Println("Generate Subset Error:", err)
		return dirs
	}
	for _, item := range items {
		dir := mk.DecodeDir(item)
		if !DirExist(dir) {
			continue
		}
		dirs = append(dirs, dir)
	}
	return dirs
}

func DecodeSets(dt Decoder, sets []string) error {
	var gw sync.WaitGroup
	for _, set := range sets {
		gw.Add(1)
		go func(set string) {
			defer gw.Done()
			if err := dt.Decode(set); err != nil {
				Err().Println("Decode set:", set, "err:", err)
			}
		}(set)
	}
	gw.Wait()
	return nil
}

func ScoreSets(ct Counter, sets []string) [][]string {
	result := [][]string{}
	for _, set := range sets {
		res, _ := ct.Score(set)
		newslice := make([][]string, len(res)+len(result))
		copy(newslice, result)
		copy(newslice[len(result):], res)
		result = newslice
	}
	return result
}

func AutoScore(identify string, dirs []string) ([][]string, error) {
	type Score struct {
		dir     string
		wer     float64
		wer_str string
	}
	rec := make(chan Score)
	go func() {
		var wg sync.WaitGroup
		for _, dir := range dirs {
			wg.Add(1)
			go func(dir string) {
				defer wg.Done()
				cmd_str := "grep WER " + dir + "/wer* | utils/best_wer.sh"
				output, err := CpuBashOutput(cmd_str)
				if err != nil || len(output) < 1 {
					Err().Println(fmt.Errorf("cmd failed or empty output!\n"))
					return
				}
				out_str := strings.Trim(string(output[:len(output)]), " \n")
				items := strings.Fields(out_str)
				if len(items) < 2 {
					Err().Println(fmt.Errorf("Output:%s  \n Output have less than 2 items\n", out_str))
					return
				}
				wer_str := items[1]
				wer, perr := strconv.ParseFloat(wer_str, 32)
				if perr != nil {
					Err().Println(perr)
					return
				}
				rec <- Score{dir, wer, wer_str}
			}(dir)
		}
		wg.Wait()
		close(rec)
	}()
	res_str := [][]string{}
	sum := 0.0
	count := 0
	// for folling the same sort
	dict := map[string]Score{}
	for sc := range rec {
		dict[sc.dir] = sc
	}

	for _, dir := range dirs {
		if sc, ok := dict[dir]; ok {
			sum += sc.wer
			count += 1
			res_str = append(res_str, []string{identify, MkScoreId(path.Base(dir)), sc.wer_str})
		}
	}
	if count > 1 { // compute the average
		res_str = append(res_str, []string{identify, MkScoreId("Ave"), fmt.Sprintf("%.2f", sum/float64(count))})
	}
	return res_str, nil

}
