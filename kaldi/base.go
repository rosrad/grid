//
package kaldi

type Model struct {
	Dst ExpBase
	Src ExpBase
	Feat
}

func NewModel() *Model {
	return &Model{*NewExpBase(), *NewExpBase(), *NewFeat()}
}
