package node

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Node interface {
	Execute() error
	GetInputs() map[string]*InputPort
	GetOutputs() map[string]*OutputPort
	Name() string
}

type InputPort struct {
	Name         string
	DefaultValue interface{}
	Connected    bool
	channel      chan interface{}
}

type OutputPort struct {
	Name     string
	channels []chan interface{}
}

type BaseNode struct {
	name    string
	inputs  map[string]*InputPort
	outputs map[string]*OutputPort
	runner  func(inputs map[string]interface{}) map[string]interface{}
	mu      sync.Mutex
}

type SubSystem struct {
	BaseNode
	system *System
}

type System struct {
	Nodes       []Node
	Connections []Connection
	mu          sync.Mutex
}

type Connection struct {
	SourceNode string
	SourcePort string
	TargetNode string
	TargetPort string
}

func NewSystem() *System {
	return &System{
		Nodes: make([]Node, 0),
	}
}

func NewBaseNode(name string) *BaseNode {
	return &BaseNode{
		name:    name,
		inputs:  make(map[string]*InputPort),
		outputs: make(map[string]*OutputPort),
	}
}

func NewSubSystem(name string) *SubSystem {
	return &SubSystem{
		BaseNode: *NewBaseNode(name),
		system:   NewSystem(),
	}
}

func (n *BaseNode) AddInputPort(name string, defaultValue interface{}) {
	n.inputs[name] = &InputPort{
		Name:         name,
		DefaultValue: defaultValue,
		channel:      make(chan interface{}, 1),
	}
}

func (n *BaseNode) AddOutputPort(name string, defaultValue interface{}) {
	n.outputs[name] = &OutputPort{
		Name:     name,
		channels: make([]chan interface{}, 0),
	}
}

func (n *BaseNode) SetRunner(runner func(inputs map[string]interface{}) map[string]interface{}) {
	n.runner = runner
}

func (n *BaseNode) Name() string {
	return n.name
}

func (n *BaseNode) GetInputs() map[string]*InputPort {
	return n.inputs
}

func (n *BaseNode) GetOutputs() map[string]*OutputPort {
	return n.outputs
}

func (sys *System) Connect(src Node, srcPort string, dest Node, destPort string) error {
	sys.mu.Lock()
	defer sys.mu.Unlock()

	// Handle subsystem port mapping
	var srcOut *OutputPort
	switch s := src.(type) {
	case *SubSystem:
		if port, ok := s.outputs[srcPort]; ok {
			srcOut = port
		} else {
			return errors.New("source port not found in subsystem")
		}
	default:
		if port, ok := src.GetOutputs()[srcPort]; ok {
			srcOut = port
		} else {
			return errors.New("source port not found")
		}
	}

	var destIn *InputPort
	switch d := dest.(type) {
	case *SubSystem:
		if port, ok := d.inputs[destPort]; ok {
			destIn = port
		} else {
			return errors.New("destination port not found in subsystem")
		}
	default:
		if port, ok := dest.GetInputs()[destPort]; ok {
			destIn = port
		} else {
			return errors.New("destination port not found")
		}
	}

	// Type checking
	srcType := reflect.TypeOf(srcOut)
	destType := reflect.TypeOf(destIn.DefaultValue)
	if reflect.TypeOf(srcOut) != reflect.TypeOf(destIn.DefaultValue) {
		return fmt.Errorf("type mismatch: %v vs %v", srcType, destType)
	}

	if destIn.Connected {
		return errors.New("destination port already connected")
	}

	srcOut.channels = append(srcOut.channels, destIn.channel)
	destIn.Connected = true

	sys.Connections = append(sys.Connections, Connection{
		SourceNode: src.Name(),
		SourcePort: srcPort,
		TargetNode: dest.Name(),
		TargetPort: destPort,
	})

	return nil
}

func (sys *System) Run() error {
	for _, node := range sys.Nodes {
		if err := node.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func (n *BaseNode) Execute() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	inputs := make(map[string]interface{})
	for name, port := range n.inputs {
		if port.Connected {
			select {
			case val := <-port.channel:
				inputs[name] = val
			default:
				if port.DefaultValue != nil {
					inputs[name] = port.DefaultValue
				}
			}
		} else if port.DefaultValue != nil {
			inputs[name] = port.DefaultValue
		}
	}

	if n.runner == nil {
		return fmt.Errorf("node %s has no runner function", n.name)
	}

	outputs := n.runner(inputs)

	for name, value := range outputs {
		if port, ok := n.outputs[name]; ok {
			for _, ch := range port.channels {
				ch <- value
			}
		}
	}
	return nil
}
