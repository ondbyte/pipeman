package tmp

import (
	"fmt"

	"github.com/ondbyte/pipeman/language_support/golang"
)

func main() {
	golang.RunCardsProgram(
		loginCard(),
	)
}

func loginCard() *golang.Card {
	return golang.NewCard("login-to-your-service").
		TakesIn("username", "").
		TakesIn("password", "").
		SpitsOut("token", "").
		Runs(
			func(in map[string]any) (out map[string]any, err error) {
				username := in["username"].(string)
				password := in["password"].(string)
				fmt.Println(username, password)
				// do any work
				out = map[string]any{}
				out["token"] = "xxxx"
				return
			},
		)
}
