//
package kaldi

import (
	"path"
	"strings"
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
	for _, set := range sets {
		if err := dt.Decode(set); err != nil {
			return err
		}
	}
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
	res_str := [][]string{}
	for _, dir := range dirs {
		cmd_str := "grep WER " + dir + "/wer* | utils/best_wer.sh |awk '{print $2}'"
		Trace().Println("Score Dir:", dir)
		output, err := BashOutput(cmd_str)
		if err != nil {
			return [][]string{}, err
		}
		wer := strings.Trim(string(output[:len(output)]), " \n")
		res_str = append(res_str, []string{identify, MkScoreId(path.Base(dir)), wer})
	}
	return res_str, nil

}
