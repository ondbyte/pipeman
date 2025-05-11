package golang

import (
	"net/http"

	"google.golang.org/protobuf/types/known/anypb"
)

type Card struct {
	name    string
	in, out map[string]any
	run     func(map[string]any) map[string]any
}

func NewCard(name string) *Card {
	return &Card{
		in:  map[string]any{},
		out: map[string]any{},
	}
}

func (c *Card) TakesIn(name string, value anypb.Any) *Card {
	var err error
	c.in[name], err = anypb.New(v)
	return c
}

func (c *Card) SpitsOut(name string, v any) *Card {
	c.out[name] = v
	return c
}

func (c *Card) Runs(run func(in map[string]any) (out map[string]any)) *Card {
	c.run = run
	return c
}

func (c *Card) attachCardToServer(s *http.ServeMux) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/in")
	endpoint := "/" + c.name
	s.Handle(endpoint, http.StripPrefix(endpoint, mux))
}
