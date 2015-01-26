//
package main

import (
	"github.com/rosrad/kaldi"
)

func main() {

	kaldi.Init()
	defer kaldi.Uninit()
	kaldi.Trace().Println("train-mono")
	m := kaldi.NewMonoTask()
	kaldi.WriteTask(m)
	if err := m.Run(); err != nil {
		kaldi.Err().Println()
	}
}
