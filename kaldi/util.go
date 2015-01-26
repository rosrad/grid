// for some utility functions used in kaldi task
package kaldi

import (
	"fmt"
	"path"
	"strings"
)

func MkGraph(target_dir string) error {
	return MkGraphOpt(target_dir, "")
}

func MkGraphOpt(target_dir, opt string) error {
	if !DirExist(target_dir) {
		Err().Println("No exist targetdir:", target_dir)
		return fmt.Errorf("MkGraph Error: No Exist TargetDir:%s \n", target_dir)
	}
	Trace().Println("mkgraph of targetdir:", target_dir)
	cmd_str := JoinArgs(
		"utils/mkgraph.sh",
		opt,
		TestLang(),
		target_dir,
		Graph(target_dir))
	Trace().Println("CMD:", cmd_str)
	LogCpuRun(cmd_str, Graph(target_dir))
	return nil
}

func PhoneSets() []string {
	return []string{"PHONESEL_dt", "PHONEMLLD_dt"}
}

func TrainName(cond string) string {
	name := "si_tr"
	if cond == "mc" {
		name = "REVERB_tr_cut/SimData_tr_for_1ch_A"
	}
	return name
}

type DataType int

const (
	ALL DataType = 1 + iota
	NOTRAIN
	DEV_TRAIN
	DEV
	DEV_TEST
	DEV_REAL
	EVAL
	REVERB_DEV
	REVERB_EVAL
	PHONE_DEV
	CLN_TRAIN
)

func DataSets(data DataType) []string {
	phone_sets := []string{"PHONE_dt", "PHONESEL_dt", "PHONEMLLD_dt"}
	dev_sets := []string{"REVERB_dt", "REVERB_CLN_dt", "REVERB_REAL_dt"}
	dev_test := []string{"REVERB_dt"}
	dev_real := []string{"REVERB_REAL_dt"}
	eval_sets := []string{"REVERB_et", "REVERB_CLN_et", "REVERB_REAL_et"}
	train_sets := []string{"REVERB_tr_cut", "si_tr"}
	cln_train_sets := []string{"si_tr"}
	switch data {
	case ALL:
		return append(append(append(phone_sets, dev_sets...), eval_sets...), train_sets...)
	case NOTRAIN:
		return append(append(phone_sets, dev_sets...), eval_sets...)
	case CLN_TRAIN:
		return cln_train_sets
	case DEV_TRAIN:
		return append(append(phone_sets, dev_sets...), train_sets...)
	case DEV_TEST:
		return dev_test
	case DEV_REAL:
		return dev_real
	case DEV:
		return append(phone_sets, dev_sets...)
	case EVAL:
		return eval_sets
	case REVERB_DEV:
		return dev_sets
	case REVERB_EVAL:
		return eval_sets
	case PHONE_DEV:
		return phone_sets
	}
	return []string{}
}

func MkScoreId(set string) string {
	return strings.Join(ExcludeDefault(Unique(set, true)), "_")
}

func FormatScore(result [][]string) string {
	if len(result) < 1 {
		Err().Println("No exist score results")
		return ""
	}
	format := ""
	split := "|-----------------------------|\n"
	format += split
	for j := 1; j < len(result[0]); j++ {
		// format += fmt.Sprintf("|%s |", result[0][0])//for Identify string
		format += fmt.Sprintf("|") //No Identify string
		for i := 0; i < len(result); i++ {

			format += fmt.Sprintf("%s|", result[i][j])
		}
		format += "\n"
	}
	format += split //
	return format
}

func FormatScoreSets(result map[string][][]string) string {
	format := ""
	split := "|-----------------------------|\n"
	for _, value := range result {
		format += split
		format += FormatScore(value)
	}
	format += split
	return format
}

func AlignPart(model, target string) string {
	return path.Join("Align", JoinDots(model, target))
}
