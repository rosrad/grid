//
package kaldi

type NormCmvn struct {
	PerUtt bool
	Mean   bool
	Vars   bool
}

func (c NormCmvn) OptStr() string {
	opt_str := ""
	// if c.Vars {
	// 	opt_str = JoinArgs(opt_str, "--cmvn-opts", `"--norm-vars=true"`)
	// }
	// if c.PerUtt {
	// 	opt_str = JoinArgs(opt_str, "--utt-cmvn", "perutt")
	// }
	return opt_str
}

func (c NormCmvn) CmdStr() string {
	key := "utt2utt"
	scp := "cmvn_utt2utt.scp"
	if !c.PerUtt {
		key = "utt2spk"
		scp = "cmvn_spk2utt.scp"
	}

	opt := ""
	if c.Vars {
		opt = JoinArgs("--norm-vars=true")
	}
	str := JoinArgs("apply-cmvn",
		opt,
		"--utt2spk=ark:"+JobStr()+"/"+key,
		"scp:"+JobStr()+"/"+scp)
	return str

}

func NewNormCmvn() *NormCmvn {
	return &NormCmvn{false, true, false}
}

func (norm Norm) NormStr() string {
	cmd := ""
	switch norm.Method {
	case "cmvn":
		cmd = norm.Cmvn.CmdStr()
	case "log":
		cmd = "apply-log"
	}
	if cmd != "" {
		return JoinArgs(cmd, "ark:-", "ark:-")
	}
	return ""
}

type Norm struct {
	Method string // cmvn,log,len
	Cmvn   *NormCmvn
}

func NewNorm() *Norm {
	return &Norm{"cmvn", NewNormCmvn()}
}
