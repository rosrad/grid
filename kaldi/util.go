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
