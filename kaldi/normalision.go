//
package kaldi

type NormCmvn struct {
	Mean bool
	Vars bool
}

func (c NormCmvn) OptStr() string {
	opt_str := ""
	if c.Vars {
		opt_str = JoinArgs(opt_str, "--cmvn-opts", `"--norm-vars=true"`)
	}
	return opt_str
}

func NewNormCmvn() *NormCmvn {
	return &NormCmvn{true, true}
}

type Norm struct {
	Cmvn *NormCmvn
}

func NewNorm() *Norm {
	return &Norm{NewNormCmvn()}
}
