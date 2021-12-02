package vcs

import (
	"scratch/pkg/scratch/app/project"
	"scratch/pkg/scratch/infrastructure/executor"
)

func NewGit(gitExecutor GitExecutor) project.VCS {
	return &gitVCS{executor: gitExecutor}
}

type gitVCS struct {
	executor GitExecutor
}

func (g *gitVCS) Init(path string) error {
	return g.executor.Run([]string{"init"}, executor.WithWorkdir(path))
}
