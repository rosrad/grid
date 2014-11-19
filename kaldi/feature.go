//
package kaldi

type FeatInterface interface {
	GetDataDir() string
	GetParamDir() string
}

type FeatBase struct {
	Dynamic string
	Norm    *Norm
}
