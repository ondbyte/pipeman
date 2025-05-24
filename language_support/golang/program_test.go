package golang_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"testing"
	"time"

	protos "github.com/ondbyte/pipeman/internal/protos/go"
	"github.com/ondbyte/pipeman/language_support/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetSocketPath() string {
	if runtime.GOOS == "windows" {
		// ensure path is absolute and short
		return os.TempDir() + "\\grpc_test.sock"
	}
	return "/tmp/grpc_test.sock"
}

func TestProg(t *testing.T) {
	r1, w2, _ := os.Pipe()
	r2, w1, _ := os.Pipe()
	go func() {
		golang.StartedByPipeman()
		err := golang.RunCardsProgram(r1, w1, "test", golang.NewCard("sum").TakesIn("a", 0).TakesIn("b", 0).SpitsOut("sum", 0).Runs(
			func(in map[string]any) (map[string]any, error) {
				return map[string]any{"sum": in["a"].(int) + in["b"].(int)}, nil
			},
		))
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(2 * time.Second)
	conn, err := grpc.NewClient("stdio", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return golang.NewStdioConn(r2, w2), nil
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := protos.NewProgramClient(conn)
	cards, err := client.GetSupportedCards(context.Background(), &protos.EmptyReq{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cards)
}
