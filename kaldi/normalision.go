//
package kaldi

type NormCmvn struct {
	PerUtt bool
	Vars   bool
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
	return &NormCmvn{false, false}
}

func (norm Norm) NormStr() string {
	cmd := ""

	if norm.Cmvn != nil {
		cmd = norm.Cmvn.CmdStr()
	} else if norm.Log != nil {
		cmd = "apply-log"
	} else if norm.Logit != nil {
		cmd = "apply-logit"
	}

	if cmd != "" {
		return JoinArgs(cmd, "ark:-", "ark:-")
	}
	return ""
}

type LogNorm struct {
}

type LogitNorm struct {
}

type Norm struct {
	Cmvn  *NormCmvn
	Log   *LogNorm
	Logit *LogitNorm
}

func NewNorm() *Norm {
	return &Norm{nil, nil, nil}
}
