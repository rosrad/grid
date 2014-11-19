package kaldi

import (
	"path"
)

func RootPath() string {
	return "tmp"
}

func Lang() string {
	return path.Join(RootPath(), "data", "lang")
}

func TestLang() string {
	return path.Join(RootPath(), "data", "lang_test_bg_5k")
}

func Graph(target string) string {
	return path.Join(target, "graph_bg_5k")
}

func JobNum(job string) int {
	parallel := 1
	switch job {
	case "decode":
		parallel = 8
	case "train":
		parallel = 8
	case "dnn":
		parallel = 16
	}
	return parallel
}
