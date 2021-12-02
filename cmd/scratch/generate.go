package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	embeddedscratch "scratch/data/scratch"
	"scratch/pkg/scratch/app/usecase"
	infragenerator "scratch/pkg/scratch/infrastructure/generator"
	"scratch/pkg/scratch/infrastructure/vcs"
)

func executeGenerate(ctx *cli.Context) error {
	loader, err := embeddedscratch.EmbedFSLoader()
	if err != nil {
		return err
	}

	scratches, err := loader.Load("manifest.yaml")
	if err != nil {
		return err
	}

	projectID := ctx.String("id")
	if projectID == "" {
		return errors.New("projectID should be not empty string")
	}

	scratchName := ctx.String("scratch")
	if projectID == "" {
		return errors.New("scratch should be not empty string")
	}

	outputPath := ctx.String("output")
	if projectID == "" {
		return errors.New("outputPath should be not empty string")
	}

	logger := initLogger()

	generator := infragenerator.NewGenerator(scratches)

	gitExecutor, err := vcs.NewGitExecutor()
	if err != nil {
		return err
	}

	u := usecase.NewGenerateProjectUseCase(generator, vcs.NewGit(gitExecutor), logger)

	return u.Execute(usecase.GenerateProjectUseCaseParams{
		ProjectID:  projectID,
		Scratch:    scratchName,
		OutputPath: outputPath,
	})
}
