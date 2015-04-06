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
	TEST DataType = 1 + iota
	DEV
	DEV_SIM
	DEV_REAL
	EVAL
	EVAL_SIM
	EVAL_REAL
	PHONE
	PHONE_DEV
	PHONE_EVAL
	PHONE_TEST
	TRAIN
	TRAIN_MC
	TRAIN_MC_DEV
	TRAIN_MC_TEST
	TRAIN_CLN
)

func SystemSet() DataType {
	return Str2DataType(SysConf().DecodeSet)
}

func Str2DataType(str string) DataType {
	type Pair struct {
		str       string
		data_type DataType
	}

	dict := []Pair{
		Pair{"TEST", TEST},
		Pair{"DEV", DEV},
		Pair{"DEV_SIM", DEV_SIM},
		Pair{"DEV_REAL", DEV_REAL},
		Pair{"EVAL", EVAL},
		Pair{"EVAL_SIM", EVAL_SIM},
		Pair{"EVAL_REAL", EVAL_REAL},
		Pair{"PHONE", PHONE},
		Pair{"PHONE_DEV", PHONE_DEV},
		Pair{"PHONE_EVAL", PHONE_EVAL},
		Pair{"PHONE_TEST", PHONE_TEST},
		Pair{"TRAIN", TRAIN},
		Pair{"TRAIN_MC", TRAIN_MC},
		Pair{"TRAIN_MC_DEV", TRAIN_MC_DEV},
		Pair{"TRAIN_MC_TEST", TRAIN_MC_TEST},
		Pair{"TRAIN_CLN", TRAIN_CLN},
	}

	key := strings.ToUpper(str)

	for _, p := range dict {
		if p.str == key {
			return p.data_type
		}
	}
	return DEV
}

func DataSets(data DataType) []string {
	phone := []string{"PHONE_dt", "PHONESEL_dt", "PHONEMLLD_dt"}

	dev_sim := []string{"REVERB_dt"}
	dev_real := []string{"REVERB_REAL_dt"}
	dev := append(dev_sim, dev_real...)

	eval_sim := []string{"REVERB_et"}
	eval_real := []string{"REVERB_REAL_et"}
	eval := append(eval_sim, eval_real...)

	test := append(dev, eval...)

	train_mc := []string{"REVERB_tr_cut"}
	train_cln := []string{"si_tr"}
	train := append(train_mc, train_cln...)

	switch data {
	case TEST:
		return append(dev, eval...)
	case DEV_SIM:
		return dev_sim
	case DEV_REAL:
		return dev_real
	case DEV:
		return eval
	case EVAL_SIM:
		return eval_sim
	case EVAL_REAL:
		return eval_real
	case EVAL:
		return eval
	case PHONE:
		return phone

	case PHONE_DEV:
		return append(dev, phone...)
	case PHONE_EVAL:
		return append(eval, phone...)
	case PHONE_TEST:
		return append(test, phone...)
	case TRAIN_MC:
		return train_mc
	case TRAIN_CLN:
		return train_cln
	case TRAIN:
		return train
	case TRAIN_MC_DEV:
		return append(train, dev...)
	case TRAIN_MC_TEST:
		return append(train, test...)
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
