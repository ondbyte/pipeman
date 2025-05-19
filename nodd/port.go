package nodd

type BasePort struct {
	node     *Node
	hasValue bool
	Val      any
	Name     string
}

func (bp *BasePort) setValue(val any) {
	bp.Val = val
	bp.hasValue = true
}

func (bp *BasePort) readValue() (any, bool) {
	return bp.Val, bp.hasValue
}

type InPort struct {
	BasePort
	Link *OutPort
}

type OutPort struct {
	BasePort
	Links []*InPort
}
