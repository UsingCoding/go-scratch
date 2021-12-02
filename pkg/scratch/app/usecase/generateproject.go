package usecase

import (
	"fmt"

	"scratch/pkg/scratch/app/project"
)

func NewGenerateProjectUseCase(
	generator project.Generator,
	vcs project.VCS,
	reporter Reporter,
) *GenerateProjectUseCase {
	return &GenerateProjectUseCase{
		generator: generator,
		vcs:       vcs,
		reporter:  reporter,
	}
}

type Reporter interface {
	Info(string ...interface{})
}

type GenerateProjectUseCase struct {
	generator project.Generator
	vcs       project.VCS
	reporter  Reporter
}

type GenerateProjectUseCaseParams struct {
	ProjectID  string
	Scratch    string
	OutputPath string
}

func (u GenerateProjectUseCase) Execute(params GenerateProjectUseCaseParams) error {
	err := u.generator.Generate(project.GenerateOpts{
		ProjectID:  params.ProjectID,
		Scratch:    params.Scratch,
		OutputPath: params.OutputPath,
	})

	u.reporter.Info(fmt.Sprintf("Project %s generated", params.ProjectID))

	if err != nil {
		return err
	}

	u.reporter.Info("Git initialized")

	return u.vcs.Init(params.OutputPath)
}
