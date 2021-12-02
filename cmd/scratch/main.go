package main

import (
	"context"
	stdlog "log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	appID = "scratch"
)

var (
	commit = "UNKNOWN"
)

func main() {
	ctx := context.Background()

	ctx = subscribeForKillSignals(ctx)

	err := runApp(ctx, os.Args)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stdlog.Fatal(string(exitErr.Stderr))
		}
		stdlog.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	app := &cli.App{
		Name: appID,
		Commands: []*cli.Command{
			{
				Name:   "generate",
				Action: executeGenerate,
				Usage:  "Generates new project from scratch",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "scratch",
						Required: true,
						Aliases:  []string{"sc"},
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Usage:    "path to output directory",
						Aliases: []string{
							"o",
						},
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Action:  executeList,
				Usage:   "Shows list of scratches",
			},
			{
				Name:   "version",
				Action: executeVersion,
				Usage:  "Prints version",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pretty",
						Aliases: []string{"p"},
					},
				},
			},
		},
	}

	return app.RunContext(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
		}
	}()

	return ctx
}
