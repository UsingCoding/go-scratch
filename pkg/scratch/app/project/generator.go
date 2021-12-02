package project

type GenerateOpts struct {
	ProjectID  string
	Scratch    string
	OutputPath string
}

type Generator interface {
	Generate(opts GenerateOpts) error
}
