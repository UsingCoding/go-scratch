package scratch

import "io/fs"

type Scratch interface {
	Name() string
	Description() string

	Structure() fs.FS
	RenameSpec() RenameSpec
}

type RenameSpec map[string]string

type scratch struct {
	name        string
	description string

	fs         fs.FS
	renameSpec RenameSpec
}

func (scratch *scratch) Name() string {
	return scratch.name
}

func (scratch *scratch) Description() string {
	return scratch.description
}

func (scratch *scratch) Structure() fs.FS {
	return scratch.fs
}

func (scratch *scratch) RenameSpec() RenameSpec {
	return scratch.renameSpec
}
