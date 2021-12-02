package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

func executeVersion(ctx *cli.Context) (err error) {
	v := versionView{
		Commit: commit,
	}

	var bytes []byte

	if ctx.Bool("pretty") {
		bytes, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
	} else {
		bytes, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	fmt.Println(string(bytes))

	return nil
}

type versionView struct {
	Commit string `json:"commit"`
}
