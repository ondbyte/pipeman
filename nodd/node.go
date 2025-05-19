package nodd

import (
	"fmt"
	"reflect"
	"sync"
)

type Runner func(map[string]any) (map[string]any, error)

type Node struct {
	Name     string
	InPorts  map[string]*InPort
	OutPorts map[string]*OutPort

	runner Runner
}

type Graph struct {
	name      string
	firstNode string
	nodes     map[string]*Node
}

func (g *Graph) SetPortValue(nodeName, portName string, val any) {
	g.nodes[nodeName].InPorts[portName].setValue(val)
}

func (g *Graph) GetFirstNode() string {
	return g.firstNode
}

func (g *Graph) SetFirstNode(name string) {
	g.firstNode = name
}

func (g *Graph) GetNodes() map[string]*Node {
	return g.nodes
}

type PortErr struct {
	Name string
	Msg  string
}
type NodeErr struct {
	Name     string
	Msg      string
	PortErrs []*PortErr
}

func (n *Node) Exec() *NodeErr {
	portErrs := []*PortErr{}
	in := map[string]any{}
	for k, v := range n.InPorts {
		if !v.hasValue {
			portErrs = append(portErrs, &PortErr{Name: k, Msg: "no value available"})
			continue
		}
		in[k] = v.Val
		v.hasValue = false
	}
	if len(portErrs) > 0 {
		return &NodeErr{
			Name:     n.Name,
			Msg:      "ports has err",
			PortErrs: portErrs,
		}
	}
	out, err := n.runner(in)
	if err != nil {
		return &NodeErr{
			Name: n.Name,
			Msg:  err.Error(),
		}
	}
	missingPortsErr := []*PortErr{}
	for k, v := range out {
		op, ok := n.OutPorts[k]
		if !ok {
			// this should never be the case
			missingPortsErr = append(missingPortsErr, &PortErr{Name: k, Msg: "the node didnt return any value for the port"})
			continue
		}
		op.Val = v
		op.hasValue = true
	}
	if len(missingPortsErr) > 0 {
		return &NodeErr{
			Name:     n.Name,
			Msg:      "ports has err",
			PortErrs: missingPortsErr,
		}
	}
	return nil
}

func exec(nodes *SyncSlice[*Node]) (*SyncSlice[*Node], *SyncSlice[*NodeErr]) {
	nextNodesToExec := NewSyncSlice[*Node]()
	execErrs := NewSyncSlice[*NodeErr]()
	wg := &sync.WaitGroup{}
	wg.Add(nodes.Len())

	//exec each node and store next linked node
	// once done with executing all nodes process err if any or move to next nodes
	for i := 0; i < nodes.Len(); i++ {
		node, _ := nodes.Get(i)
		go func() {
			defer wg.Done()
			err := node.Exec()
			if err != nil {
				execErrs.Append(err)
				return
			}
			portErrs := []*PortErr{}
			for _, p := range node.OutPorts {
				if !p.hasValue {
					portErrs = append(portErrs, &PortErr{Name: p.Name, Msg: "port has no value set"})
					continue
				}
				for _, l := range p.Links {
					l.Val = p.Val
					nextNodesToExec.Append(l.node)
				}
				p.hasValue = false
			}
			if len(portErrs) > 0 {
				execErrs.Append(
					&NodeErr{
						Name:     node.Name,
						Msg:      "ports has err",
						PortErrs: portErrs,
					},
				)
			}
		}()
	}
	wg.Wait()
	if execErrs.Len() > 0 {
		return nil, execErrs
	}
	return nextNodesToExec, nil
}

type GraphErr struct {
	Name     string
	Msg      string
	NodeErrs []*NodeErr
}

func (g *Graph) Exec() *GraphErr {
	if g.firstNode == "" {
		return &GraphErr{
			Name:     g.name,
			Msg:      "no first node set",
			NodeErrs: nil,
		}
	}
	nodes := NewSyncSlice[*Node]()
	nodes.Append(g.nodes[g.firstNode])
	for {
		var errs *SyncSlice[*NodeErr]
		nodes, errs = exec(nodes)
		if errs.Len() > 0 {
			return &GraphErr{
				Name:     g.name,
				Msg:      "nodes has err",
				NodeErrs: errs.Slice(),
			}
		}
		if nodes.Len() == 0 {
			break
		}
	}
	return nil
}

func (g *Graph) AddNode(node *Node) {
	g.nodes[node.Name] = node
}

func (g *Graph) CanLinkPort(outNode, outPort, InNode, InPort string) error {
	return g.linkPort(outNode, outPort, InNode, InPort, false)
}

func (g *Graph) LinkPort(outNode, outPort, InNode, InPort string) error {
	return g.linkPort(outNode, outPort, InNode, InPort, true)
}

func (g *Graph) linkPort(outNode, outPort, InNode, InPort string, link bool) error {
	na, ok := g.nodes[outNode]
	if !ok {
		return fmt.Errorf("node '%v' doesnt exist", outNode)
	}
	op, ok := na.OutPorts[outPort]
	if !ok {
		return fmt.Errorf("port '%v' doesnt exist", outPort)
	}
	nb, ok := g.nodes[InNode]
	if !ok {
		return fmt.Errorf("node '%v' doesnt exist", InNode)
	}
	ip, ok := nb.InPorts[InPort]
	if !ok {
		return fmt.Errorf("port '%v' doesnt exist", InPort)
	}
	//check type
	if reflect.TypeOf(op.Val) != reflect.TypeOf(ip.Val) {
		return fmt.Errorf("type mismatch")
	}
	if !link {
		return nil
	}
	op.Links = append(op.Links, ip)
	ip.Link = op
	return nil
}

func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]*Node{},
	}
}
