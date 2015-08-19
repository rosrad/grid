package kaldi

import (
	"net"
	"net/rpc"
	"time"
)

type NodeRPC struct {
	DV *DevSel
}

func (n *NodeRPC) AutoSelectGpu(key string, rn *Node) error {
	*rn = n.DV.AutoSelectGpu()
	return nil
}

func (n *NodeRPC) AutoSelectCpu(key string, rn *Node) error {
	*rn = n.DV.AutoSelectCpu()
	return nil
}

func (n *NodeRPC) Update(key string, rn *bool) error {
	n.DV.Update()
	*rn = true
	return nil
}

func (n *NodeRPC) Inited(key string, rn *bool) error {
	*rn = n.DV.Inited()
	return nil
}

func (n *NodeRPC) GpuSort(key string, ns *[]Node) error {
	n.DV.SortGpu()
	for _, n := range n.DV.Nodes {
		*ns = append(*ns, *n)
	}
	return nil
}

func (n *NodeRPC) CpuSort(key string, ns *[]Node) error {
	n.DV.SortCpu()
	for _, n := range n.DV.Nodes {
		*ns = append(*ns, *n)
	}
	return nil
}

type RPCClient struct {
	conn *rpc.Client
}

func NewClient(ip string) *RPCClient {
	conn, err := net.DialTimeout("tcp", ip, time.Millisecond*500)
	if err != nil {
		Err().Println("dialing:", err)
		return &RPCClient{}
	}

	return &RPCClient{conn: rpc.NewClient(conn)}
}

func (c *RPCClient) Inited() bool {

	if c.conn == nil {
		return false
	}
	inited := false
	err := c.conn.Call("NodeRPC.Inited", "", &inited)
	if err != nil {
		Err().Println("Inited error:", err)
	}
	return inited
}

func (c *RPCClient) AutoSelectGpu() Node {
	nd := Node{}
	err := c.conn.Call("NodeRPC.AutoSelectGpu", "", &nd)
	if err != nil {
		Err().Println("arith error:", err)
	}
	return nd
}

func (c *RPCClient) AutoSelectCpu() Node {
	nd := Node{}
	err := c.conn.Call("NodeRPC.AutoSelectCpu", "", &nd)
	if err != nil {
		Err().Println("arith error:", err)
	}
	return nd
}

func (c *RPCClient) GpuSort() []Node {
	nd := []Node{}
	err := c.conn.Call("NodeRPC.GpuSort", "", &nd)
	if err != nil {
		Err().Println("arith error:", err)
	}
	return nd
}

func (c *RPCClient) Update() {
	res := false
	err := c.conn.Call("NodeRPC.Update", "", &res)
	if err != nil {
		Err().Println("arith error:", err)
	}
}
func (c *RPCClient) CpuSort() []Node {
	nd := []Node{}
	err := c.conn.Call("NodeRPC.CpuSort", "", &nd)
	if err != nil {
		Err().Println("arith error:", err)
	}
	return nd
}
