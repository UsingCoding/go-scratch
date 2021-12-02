package vcs

import "scratch/pkg/scratch/infrastructure/executor"

const git = "git"

type GitExecutor = executor.Executor

func NewGitExecutor() (GitExecutor, error) {
	return executor.New(git)
}
