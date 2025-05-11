package main

import "github.com/urfave/cli/v3"

func main() {
	cmd := cli.Command{
		Name: "pman",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "test",
			},
		},
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}
