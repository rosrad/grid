// control the basic path structure
package kaldi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ExpBase struct {
	Name  string
	Exp   string
	Feat  string
	Label string
}

func NewExpBase() *ExpBase {
	return &ExpBase{Name: "tri1", Exp: "GMM", Feat: "mfcc", Label: "normal"}
}

func (p ExpBase) ExpDir() string {
	return path.Join(RootPath(), "exp", p.Feat, p.Label, p.Exp, p.Name)
}
func (p ExpBase) DataDir() string {
	return path.Join(RootPath(), "feats", p.Feat, p.Label, "data")
}

func (p ExpBase) ParamDir() string {
	return path.Join(RootPath(), "feats", p.Feat, p.Label, "data")
}
func (p ExpBase) LogDir() string {
	return path.Join(RootPath(), "log", p.Feat, p.Label, p.Exp)
}

func (p ExpBase) DeriveExp() string {
	return path.Join(RootPath(), "feats", p.Feat, "exp", p.Name)
}

func (p ExpBase) DeriveDump() string {
	return path.Join(RootPath(), "feats", p.Feat, "exp", p.Name, "dump")
}

func (p ExpBase) AlignDir() string {
	a := p
	a.Name = JoinDots(a.Exp, a.Name)
	a.Exp = "ALIGN"
	return a.ExpDir()
}

func (p ExpBase) TrainData(cond string) string {
	train := "si_tr"
	if cond == "mc" {
		train = "REVERB_tr_cut/SimData_tr_for_1ch_A"
	}
	return path.Join(p.DataDir(), train)
}

func (p ExpBase) StoreName() map[string]string {
	names := make(map[string]string)
	names["exp"], _ = filepath.EvalSymlinks(p.ExpDir())
	names["data"], _ = filepath.EvalSymlinks(p.DataDir())
	names["param"], _ = filepath.EvalSymlinks(p.ParamDir())
	return names
}

func (p ExpBase) Bakup() int {
	num := 0

	timestamp := time.Now().Format("20060102150405")
	for _, value := range p.StoreName() {
		if !DirExist(value) {
			bakup := value + ".bakup." + timestamp
			Trace().Printf("Mv %s, %s \n", value, bakup)
			os.Rename(value, bakup)
			num++
		}
	}
	return num
}

func (p ExpBase) Subsets(set string) ([]string, error) {
	sets := []string{set}
	if !strings.HasPrefix(set, "si_") {
		sys_dirs, err := ioutil.ReadDir(path.Join(p.DataDir(), set))
		if err != nil || len(sys_dirs) == 0 {
			return []string{}, fmt.Errorf("No Subset in %s", set)
		}
		sets = []string{}
		for _, tmp := range sys_dirs {
			sets = append(sets, path.Join(set, tmp.Name()))
		}
	}
	return sets, nil
}