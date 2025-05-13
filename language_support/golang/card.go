package golang

type Card struct {
	name    string
	in, out map[string]any
	run     func(map[string]any) (map[string]any, error)
}

func NewCard(name string) *Card {
	return &Card{
		in:  map[string]any{},
		out: map[string]any{},
	}
}

// arg type/default values this card takes
func (c *Card) TakesIn(name string, value any) *Card {
	c.in[name] = value
	return c
}

// values which this card spits out
func (c *Card) SpitsOut(name string, v any) *Card {
	c.out[name] = v
	return c
}

// passed fn which runs when this card is called with arg 'in' which you defined using 'TakesIn' method and
// returns the values you defined using 'SpitsOut' method
func (c *Card) Runs(run func(in map[string]any) (map[string]any, error)) *Card {
	c.run = run
	return c
}
