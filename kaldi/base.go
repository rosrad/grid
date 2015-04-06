//
package kaldi

type Model struct {
	Dst ExpBase
	Src ExpBase
	Ali ExpBase
	Feat
}

func NewModel() *Model {
	return &Model{*NewExpBase(), *NewExpBase(), *NewExpBase(), *NewFeat()}
}
