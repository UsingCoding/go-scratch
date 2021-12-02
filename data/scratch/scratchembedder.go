package scratch

import (
	"embed"
	"io/fs"

	"scratch/pkg/scratch/infrastructure/scratch"
)

// Ignored files '.' excluded from embedding, will be fixed with
// https://github.com/golang/go/issues/43854

//go:embed templates/*
var scratches embed.FS

func EmbedFSLoader() (scratch.Loader, error) {
	sub, err := fs.Sub(scratches, "templates")
	if err != nil {
		return nil, err
	}

	return scratch.NewFsLoader(sub), nil
}
