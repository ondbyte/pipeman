package golang

import (
	"context"
	"fmt"
	"net"
	"os"

	protos "github.com/ondbyte/pipeman/internal/protos/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type program struct {
	cards []*Card
	protos.UnimplementedProgramServer
}

// GetSupportedCards implements _go.ProgramServer.
func (p *program) GetSupportedCards(context.Context, *protos.EmptyReq) (cards *protos.Cards, err error) {
	cards = &protos.Cards{
		Cards: []*protos.Card{},
	}
	for _, card := range p.cards {
		cards.Cards = append(cards.Cards, &protos.Card{
			Name: card.name,
			In: &protos.CardInput{
				Entries: map[string]*anypb.Any{
					"username": anypb.New(0),
				},
			},
			Out: card.Output,
		})
	}
}

// RunCard implements _go.ProgramServer.
func (p *program) RunCard(context.Context, *protos.CardInputWithCardName) (*protos.CardOutput, error) {
	panic("unimplemented")
}

var _ protos.ProgramServer = (*program)(nil)

func RunCardsProgram(cards ...*Card) error {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return fmt.Errorf("PORT environment variable not set")
	}
	s := grpc.NewServer()
	protos.RegisterProgramServer(s, &program{cards: cards})

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
