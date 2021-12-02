package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/liggitt/tabwriter"

	embeddedscratch "scratch/data/scratch"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = tabwriter.RememberWidths
)

func executeList(ctx *cli.Context) error {
	loader, err := embeddedscratch.EmbedFSLoader()
	if err != nil {
		return err
	}

	scratches, err := loader.Load("manifest.yaml")
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)

	_, err = fmt.Fprintln(w, "Name\tDescription")
	if err != nil {
		return err
	}

	for _, scratch := range scratches {
		_, err = fmt.Fprintf(w, "%s\t%s\n", scratch.Name(), scratch.Description())
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
