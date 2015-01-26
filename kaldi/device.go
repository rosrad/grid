//

package kaldi

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Node struct {
	Node    string
	GpuMem  float32
	LoadAve float32
	CpuNum  int
	CpuMHz  float32
}

func NewNode() *Node {
	return &Node{"", 0, 0.0, 0, 0.0}
}
func (gn Node) GpuUsage() float32 {
	return gn.GpuMem + gn.CpuUsage()/100
}

func (gn Node) CpuUsage() float32 {
	return float32(gn.CpuNum) * gn.CpuMHz * (100 - gn.LoadAve) / 100
}

func (gn *Node) SyncInfo() error {
	cmd_str := JoinArgs(
		"ssh ", gn.Node, "bash -c",
		"sys-info.sh")
	out, err := BashOutput(cmd_str)
	if err != nil {
		return err
	}
	str := strings.Trim(string(out), "\n ")
	for idx, value := range strings.Split(str, "\n") {
		switch idx {
		case 0:
			load, err1 := strconv.ParseFloat(value, 32)
			gn.LoadAve = float32(load)
			if err1 != nil {
				err = err1
			}
		case 1:
			mem, err2 := strconv.ParseFloat(value, 32)
			gn.GpuMem = float32(mem)
			if err2 != nil {
				err = err2
			}
		case 2:
			num, err3 := strconv.ParseInt(value, 0, 32)
			gn.CpuNum = int(num)
			if err3 != nil {
				err = err3
			}
		case 3:
			cpu_freq, err4 := strconv.ParseFloat(value, 32)
			gn.CpuMHz = float32(cpu_freq)
			if err4 != nil {
				err = err4
			}
		}
	}
	return err
}

type DevSel struct {
	Nodes []*Node
	m     int
	quit  <-chan struct{}
	msg   chan int
}

func NewDevSel() *DevSel {
	return &DevSel{m: 60}
}

var dev_instance *DevSel

func DevInstance() *DevSel {
	if dev_instance == nil {
		dev_instance = NewDevSel()
	}
	return dev_instance
}

func (gs *DevSel) AutoSync() {
	dev_instance.Init()
	log.Println("DevSel AutoSync Start...")
	dev_instance.SyncInfo()
	dev_instance.autoSync()
}

func contains(s []int, d int) bool {
	for _, i := range s {
		if d == i {
			return true
		}
	}
	return false
}

func (gs *DevSel) Init() {
	exclude := []int{8, 11}
	const MaxNode = 13
	for i := 1; i < MaxNode+1; i++ {
		if contains(exclude, i) {
			fmt.Println("Exclude", i)
			continue
		}
		n := NewNode()
		n.Node = fmt.Sprintf("node%02d", i)
		gs.Nodes = append(gs.Nodes, n)
	}
}

func (gs *DevSel) autoSync() {
	m := gs.m
	if m <= 0 {
		m = 5
	}
	ticker := time.NewTicker(time.Duration(m) * time.Minute)
	gs.quit = make(chan struct{})
	gs.msg = make(chan int, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				gs.SyncInfo()
			case <-gs.msg:
				gs.SyncInfo()
			case <-gs.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (gs *DevSel) SyncInfo() {
	log.Println("")
	log.Println("Geting system info...")
	var wg sync.WaitGroup
	for _, v := range gs.Nodes {
		wg.Add(1)
		sync := func(gn *Node) {
			defer wg.Done()
			if err := gn.SyncInfo(); err != nil {
				log.Println("Node :", gn.Node)
				log.Println("SyncInfo", err)
			}
		}
		go sync(v)
	}
	wg.Wait()
	log.Println("Obtained system info!")
	log.Println("")
}

func (gs *DevSel) SortGpu() {
	mem := func(n1, n2 *Node) bool {
		return n1.GpuUsage() > n2.GpuUsage()
	}
	By(mem).Sort(gs.Nodes)
}

func (gs *DevSel) SortCpu() {
	load := func(n1, n2 *Node) bool {
		return n1.CpuUsage() > n2.CpuUsage()
	}
	By(load).Sort(gs.Nodes)
}

func (gs *DevSel) update() {
	if gs.msg != nil {
		select {
		case msg, ok := <-gs.msg:
			if ok {
				fmt.Println("")
				fmt.Println("getting system info with msg", msg)
				fmt.Println("")
				gs.msg <- msg + 1
			}
		default:
			gs.msg <- 1
		}
	}
}
func (gs *DevSel) AutoSelectGpu() Node {
	gs.SortGpu()
	opt_node := gs.Nodes[0]
	gs.Nodes[0].GpuMem = gs.Nodes[0].GpuMem / 2
	gs.update()
	return *opt_node
}

func (gs *DevSel) AutoSelectCpu() Node {
	gs.SortCpu()
	opt_node := gs.Nodes[0]
	gs.Nodes[0].LoadAve += 30
	gs.update()
	return *opt_node
}

type By func(p1, p2 *Node) bool
type NodeSorter struct {
	nodes []*Node
	by    By
}

// Len is part of sort.Interface.
func (s *NodeSorter) Len() int {
	return len(s.nodes)
}

// Swap is part of sort.Interface.
func (s *NodeSorter) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *NodeSorter) Less(i, j int) bool {
	return s.by(s.nodes[i], s.nodes[j])
}

func (by By) Sort(nodes []*Node) {
	ps := &NodeSorter{
		nodes: nodes,
		by:    by}
	sort.Sort(ps)
}

func (gs DevSel) PrintNodes(bGpu bool) {
	for _, n := range gs.Nodes {
		if bGpu {
			log.Printf("Node:%s, GpuUsage:%04.2f, GpuMem:%04.2f", n.Node, n.GpuUsage(), n.GpuMem)
		} else {
			log.Printf("Node:%s, CpuUsage:%04.2f, CPU(s):%02d, CPU MHz: %04.2f, LoadAve:%02.2f\n",
				n.Node, n.CpuUsage(), n.CpuNum, n.CpuMHz, n.LoadAve)
		}
	}
}
