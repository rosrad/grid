// control the basic path structure
package kaldi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ExpBase struct {
	Name  string
	Exp   string
	Feat  string
	Label string
	Root  string
}

func NewExpBase() *ExpBase {
	return &ExpBase{Name: "tri1", Exp: "GMM", Feat: "mfcc", Label: "normal", Root: ""}
}

func (p ExpBase) RootPath() string {
	if p.Root == "" {
		return RootPath()
	}
	return p.Root
}

func (p ExpBase) ExpDir() string {
	return path.Join(p.RootPath(), "exp", p.Feat, p.Label, p.Exp, p.Name)
}
func (p ExpBase) DataDir() string {
	return path.Join(p.RootPath(), "feats", p.Feat, p.Label, "data")
}

func (p ExpBase) ParamDir() string {
	return path.Join(p.RootPath(), "feats", p.Feat, p.Label, "param")
}
func (p ExpBase) LogDir() string {
	return path.Join(p.RootPath(), "log", p.Feat, p.Label, p.Exp)
}

func (p ExpBase) DeriveExp() string {
	return path.Join(p.RootPath(), "feats", p.Feat, "exp", p.Label)
}

func (p ExpBase) DeriveDump() string {
	return path.Join(p.RootPath(), "feats", p.Feat, "exp", p.Label, "dump")
}

func (p ExpBase) AlignDir() string {
	a := p
	a.Name = JoinDots(a.Exp, a.Name)
	a.Exp = "ALIGN"
	return a.ExpDir()
}

func (p ExpBase) TrainData(cond string) string {
	return path.Join(p.DataDir(), TrainName(cond))
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
	sets := []string{}
	if FileExist(path.Join(p.DataDir(), set, "feats.scp")) == nil {
		return []string{set}, nil
	}

	sys_dirs, _ := ioutil.ReadDir(path.Join(p.DataDir(), set))
	for _, tmp := range sys_dirs {
		if !tmp.IsDir() ||
			FileExist(path.Join(p.DataDir(), set, tmp.Name(), "feats.scp")) != nil {
			continue
		}
		sets = append(sets, path.Join(set, tmp.Name()))
	}

	if len(sets) == 0 {
		return []string{}, fmt.Errorf("No Exist dataset")
	}
	return sets, nil
}
