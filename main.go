package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ahmedtofaha10/doc-t/laravel"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Usage:                "make documentations for your projects ;)",
		Commands: []*cli.Command{
			{
				Name:    "laravel",
				Aliases: []string{"L"},
				Usage:   "start document of laravel project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Usage:    "path for the base/root directory of the project",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					path := cCtx.String("path")
					laravel.Documenting(path)
					fmt.Println("start documenting:", path)
					return nil
				},
				BashComplete: func(cCtx *cli.Context) {
					// This will complete if no args are passed
					if cCtx.NArg() > 0 {
						return
					}
					fmt.Println("TASK NAME")
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
