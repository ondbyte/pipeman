package node

import (
	"testing"
)

func TestBasicFlow(t *testing.T) {
	sys := NewSystem()

	// Create nodes
	producer := NewBaseNode("producer")
	processor := NewBaseNode("processor")
	consumer := NewBaseNode("consumer")

	// Add ports
	producer.AddOutputPort("numbers", 0)
	processor.AddInputPort("in", 0)
	processor.AddOutputPort("out", 0)
	consumer.AddInputPort("result", 0)

	// Configure runners
	producer.SetRunner(func(inputs map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"numbers": 21}
	})

	processor.SetRunner(func(inputs map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"out": inputs["in"].(int) * 2}
	})

	consumer.SetRunner(func(inputs map[string]interface{}) map[string]interface{} {
		if inputs["result"].(int) != 42 {
			t.Errorf("Expected 42, got %d", inputs["result"])
		}
		return nil
	})

	// Add nodes to system
	sys.Nodes = append(sys.Nodes, producer, processor, consumer)

	// Connect nodes
	if err := sys.Connect(producer, "numbers", processor, "in"); err != nil {
		t.Fatal(err)
	}
	if err := sys.Connect(processor, "out", consumer, "result"); err != nil {
		t.Fatal(err)
	}

	// Execute system
	if err := sys.Run(); err != nil {
		t.Fatal(err)
	}
}
func TestSubsystem(t *testing.T) {
	sys := NewSystem()

	// Create main nodes
	producer := NewBaseNode("producer")
	producer.AddOutputPort("out", 0)
	producer.SetRunner(func(inputs map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"out": 21}
	})

	// Create subsystem
	sub := NewSubSystem("multiplier")
	sub.AddInputPort("in", 0)
	sub.AddOutputPort("out", 0)

	// Internal multiplier node
	multiplier := NewBaseNode("multiplier")
	multiplier.AddInputPort("in", 0)
	multiplier.AddOutputPort("out", 0)
	multiplier.SetRunner(func(inputs map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"out": inputs["in"].(int) * 2}
	})

	// Connect internal nodes
	sub.system.Nodes = append(sub.system.Nodes, multiplier)
	sub.system.Connect(sub, "in", multiplier, "in")
	sub.system.Connect(multiplier, "out", sub, "out")

	// Create consumer
	consumer := NewBaseNode("consumer")
	consumer.AddInputPort("in", 0)
	consumer.SetRunner(func(inputs map[string]interface{}) map[string]interface{} {
		if inputs["in"].(int) != 42 {
			t.Errorf("Expected 42, got %d", inputs["in"])
		}
		return nil
	})

	// Build system
	sys.Nodes = append(sys.Nodes, producer, sub, consumer)

	// Connect main system
	if err := sys.Connect(producer, "out", sub, "in"); err != nil {
		t.Fatal(err)
	}
	if err := sys.Connect(sub, "out", consumer, "in"); err != nil {
		t.Fatal(err)
	}

	// Execute and verify
	if err := sys.Run(); err != nil {
		t.Fatalf("Execution failed: %v", err)
	}
}
func TestTypeValidation(t *testing.T) {
	sys := NewSystem()

	n1 := NewBaseNode("n1")
	n2 := NewBaseNode("n2")

	n1.AddOutputPort("out", "string")
	n2.AddInputPort("in", 0)

	err := sys.Connect(n1, "out", n2, "in")
	if err == nil {
		t.Fatal("Expected type mismatch error")
	}

	expectedErr := "type mismatch: string vs int"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
