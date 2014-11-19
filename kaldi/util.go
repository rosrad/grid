// for some utility functions used in kaldi task
package kaldi

import (
	"fmt"
	"path"
	"strings"
)

func MkGraph(target_dir string) error {
	if !DirExist(target_dir) {
		Err().Println("No exist targetdir:", target_dir)
		return fmt.Errorf("MkGraph Error: No Exist TargetDir:%s \n", target_dir)
	}
	Trace().Println("mkgraph of targetdir:", target_dir)
	cmd_str := JoinArgs(
		"utils/mkgraph.sh",
		TestLang(),
		target_dir,
		Graph(target_dir))
	Trace().Println("CMD:", cmd_str)
	BashRun(cmd_str)
	return nil
}

func PhoneSets() []string {
	return []string{"PHONESEL_dt", "PHONEMLLD_dt"}
}

func DataSets(dt_only bool) []string {
	if dt_only {
		return []string{"REVERB_dt", "PHONE_dt", "PHONESEL_dt", "PHONEMLLD_dt"}
	}
	return []string{"REVERB_tr_cut", "REVERB_dt", "PHONE_dt", "PHONESEL_dt", "PHONEMLLD_dt"}
}

func MkScoreId(set string) string {
	return strings.Join(ExcludeDefault(Unique(set, true)), "_")
}

func FormatScore(result [][]string) string {
	format := ""
	split := "|-----------------------------|\n"
	format += split
	for j := 1; j < len(result[0]); j++ {
		format += fmt.Sprintf("|%s |", result[0][0])
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
