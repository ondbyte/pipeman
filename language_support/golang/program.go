package golang

import (
	"context"
	"fmt"
	"os"

	protos "github.com/ondbyte/pipeman/internal/protos/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

// program is collection of cards
type program struct {
	cards map[string]*Card
	protos.UnimplementedProgramServer
}

// GetSupportedCards implements _go.ProgramServer.
func (p *program) GetSupportedCards(context.Context, *protos.EmptyReq) (pCards *protos.Cards, err error) {
	pCards = &protos.Cards{
		Cards: []*protos.Card{},
	}
	for _, c := range p.cards {
		in, err := structpb.NewStruct(c.in)
		if err != nil {
			return nil, fmt.Errorf("failed to convert input to struct: %v", err)
		}
		out, err := structpb.NewStruct(c.out)
		if err != nil {
			return nil, fmt.Errorf("failed to convert output to struct: %v", err)
		}
		pCards.Cards = append(pCards.Cards, &protos.Card{
			Name: c.name,
			In:   in,
			Out:  out,
		})
	}
	return pCards, nil
}

// RunCard implements _go.ProgramServer.
func (p *program) RunCard(ctx context.Context, iwcn *protos.CardInputWithCardName) (*structpb.Struct, error) {
	in := iwcn.Input.AsMap()
	card, ok := p.cards[iwcn.Card]
	if !ok {
		return nil, fmt.Errorf("card %s not found", iwcn.Card)
	}
	m, err := card.run(in)
	if err != nil {
		return nil, fmt.Errorf("failed to run card: %v", err)
	}
	r, err := structpb.NewStruct(m)
	if err != nil {
		return nil, fmt.Errorf("failed to convert map to struct: %v", err)
	}
	return r, nil
}

var _ protos.ProgramServer = (*program)(nil)

func newProgram(cards ...*Card) *program {
	prgm := &program{cards: map[string]*Card{}}
	for _, c := range cards {
		prgm.cards[c.name] = c
	}
	return prgm
}

type ProgramData struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

var (
	startedByKey         = "STARTED_BY"
	startedByVal         = "pipeman-random-yadhu"
	StartedByEnvKeyValue = fmt.Sprintf("%s=%s", startedByKey, startedByVal)
)

func StartedByPipeman() {
	err := os.Setenv(startedByKey, startedByVal)
	if err != nil {
		panic(err)
	}
}

// takes your cards and runs them as single go program,
// you should call this from your main, no other requirements
func RunCardsProgram(r, w *os.File, programName string, cards ...*Card) error {
	if os.Getenv(startedByKey) != startedByVal {
		// not started by pipeman
		return fmt.Errorf("exiting because not started by pipeman")
	}
	s := grpc.NewServer()
	protos.RegisterProgramServer(s, newProgram(cards...))

	lis := &StdioListener{
		r: r, w: w,
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
