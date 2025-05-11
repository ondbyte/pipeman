package tmp

import (
	"fmt"

	"github.com/ondbyte/pipeman/card"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	card.RunCardsProgram(
		loginCard(),
	)
}

func loginCard() *card.Card {
	return card.NewCard().
		TakesIn("username", wrapperspb.String("")).
		TakesIn("password", "").
		SpitsOut("token", "").
		Runs(
			func(in map[string]any) (out map[string]any) {
				username := in["username"].(*wrapperspb.StringValue)
				password := in["password"]
				fmt.Println(username, password)
				// do any work
				out = map[string]any{}
				out["token"] = "xxxx"
				return
			},
		)

}
