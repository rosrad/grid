//
package main

import (
	"github.com/rosrad/kaldi"
	"log"
	"time"
)

func main() {
	// gs := kaldi.NewDevSel()
	// gs.Init()
	// gs.SyncInfo()
	// gs.SortGpu()
	// gs.PrintNodes()
	kaldi.DevInstance().AutoSync()
	for {
		kaldi.DevInstance().SortGpu()
		kaldi.DevInstance().PrintNodes(true)
		log.Println("")
		kaldi.DevInstance().SortCpu()
		kaldi.DevInstance().PrintNodes(false)
		log.Println("")
		kaldi.DevInstance().AutoSelectGpu()
		time.Sleep(15 * time.Second)
	}
}
