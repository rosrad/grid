package kaldi

import (
	"encoding/json"
	"os"
	"strings"
)

type DatasetConf struct {
	CLN_tr string
	MC_tr  string
	Dev    []string
	Eval   []string
}

var dataset_conf = &DatasetConf{"", "", []string{}, []string{}}

func LoadDataConf() error {
	data_conf := "dataset.json"
	f, err := os.Open(data_conf)
	if err != nil {
		Err().Println("Dataset Conf Read :", err)
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(Dataset())
	if err != nil {
		Err().Println("Dataset Conf Decode :", err)
		return err
	}
	return nil
}

func Dataset() *DatasetConf {
	return dataset_conf
}

func TrainName(cond string) string {
	name := Dataset().CLN_tr
	// name = "REVERB_tr_cut/SimData_tr_for_1ch_A"
	if cond == "mc" {
		name = Dataset().MC_tr
	}
	return name
}

type DataType int

const (
	TEST DataType = 1 + iota
	DEV
	EVAL
	TRAIN
	TRAIN_CLN
	TRAIN_MC
	TRAIN_MC_DEV
	TRAIN_MC_TEST
)

func DataSets(data DataType) []string {
	dev := Dataset().Dev
	eval := Dataset().Eval
	test := append(dev, eval...)
	train_mc := []string{Dataset().MC_tr}
	train_cln := []string{Dataset().CLN_tr}
	train := append(train_mc, train_cln...)

	switch data {
	case TEST:
		return append(dev, eval...)
	case DEV:
		return dev
	case EVAL:
		return eval
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
		Pair{"EVAL", EVAL},
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
