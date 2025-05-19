package nodd

import "testing"

func TestNodeFlow(t *testing.T) {
	sumInputPortA := &InPort{
		BasePort: BasePort{
			Name: "a",
			Val:  0,
		},
	}
	sumInputPortB := &InPort{
		BasePort: BasePort{
			Name: "b",
			Val:  0,
		},
	}
	sumOutPort := &OutPort{
		BasePort: BasePort{
			Name: "sum",
			Val:  0,
		},
	}
	sumNode := &Node{
		Name: "sum",
		InPorts: map[string]*InPort{
			sumInputPortA.Name: sumInputPortA,
			sumInputPortB.Name: sumInputPortB,
		},
		OutPorts: map[string]*OutPort{
			sumOutPort.Name: sumOutPort,
		},
		runner: func(m map[string]any) (r map[string]any, err error) {
			r = map[string]any{
				sumOutPort.Name: m[sumInputPortA.Name].(int) + m[sumInputPortB.Name].(int),
			}
			return
		},
	}
	productInputA := &InPort{
		BasePort: BasePort{
			Name: "a",
			Val:  0,
		},
	}

	productInputB := &InPort{
		BasePort: BasePort{
			Name: "b",
			Val:  0,
		},
	}
	productOutput := &OutPort{
		BasePort: BasePort{
			Name: "product",
			Val:  0,
		},
	}
	productNode := &Node{
		Name: "product",
		InPorts: map[string]*InPort{
			productInputA.Name: productInputA,
			productInputB.Name: productInputB,
		},
		OutPorts: map[string]*OutPort{
			productOutput.Name: productOutput,
		},
		runner: func(m map[string]any) (r map[string]any, err error) {
			r = map[string]any{
				"product": m[sumOutPort.Name].(int) * 2,
			}
			return
		},
	}

	graph := NewGraph()
	graph.AddNode(sumNode)
	graph.AddNode(productNode)
	graph.SetPortValue(sumOutPort.Name, sumInputPortA.Name, 100)
	graph.SetPortValue(sumOutPort.Name, sumInputPortB.Name, 200)
	graph.LinkPort(sumNode.Name, sumOutPort.Name, productNode.Name, productOutput.Name)
	graph.SetPortValue(productOutput.Name, productInputB.Name, 2)
	graph.SetFirstNode(sumNode.Name)
	err := graph.Exec()
	if err != nil {
		t.Error(err)
	}
}
