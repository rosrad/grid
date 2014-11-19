//
package kaldi

type Model struct {
	Src ExpBase
	Dst ExpBase
}

func NewModel() *Model {
	return &Model{*NewExpBase(), *NewExpBase()}
}

type ModelConf struct {
	Dynamic string // raw,delta
	Norm    *Norm  // Normalision config
}

func NewModelConf() *ModelConf {
	return &ModelConf{"raw", NewNorm()}
}
func (conf ModelConf) OptStr() string {
	opt := conf.Norm.Cmvn.OptStr()
	if conf.Dynamic != "" {
		opt = JoinArgs(opt, "--feat-type", conf.Dynamic)
	}
	return opt
}
