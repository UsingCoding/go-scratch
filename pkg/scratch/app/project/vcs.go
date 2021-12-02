package project

type VCS interface {
	Init(path string) error
}
