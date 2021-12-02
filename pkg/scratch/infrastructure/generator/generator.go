package generator

import (
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/pkg/errors"

	"scratch/pkg/scratch/app/project"
	"scratch/pkg/scratch/infrastructure/scratch"
)

func NewGenerator(scratches []scratch.Scratch) project.Generator {
	return &generator{scratches: scratches}
}

type generator struct {
	scratches []scratch.Scratch
}

func (generator *generator) Generate(opts project.GenerateOpts) error {
	if stat, err := os.Stat(opts.OutputPath); os.IsNotExist(err) || !stat.IsDir() {
		return errors.Errorf("outputh path is not exists or it`s not a directory %s", opts.OutputPath)
	}

	s, err := generator.findScratch(opts.Scratch)
	if err != nil {
		return err
	}

	err = fs.WalkDir(s.Structure(), ".", func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entryPath == "." {
			return nil
		}

		outPath := path.Join(opts.OutputPath, entryPath)

		outPath = replaceTemplateVariablesInPath(outPath, opts.ProjectID)

		outPath = generator.updateOutPathByScratchRenameSpec(outPath, s)

		if d.IsDir() {
			if dirInfo, err2 := os.Stat(outPath); err2 == nil {
				if dirInfo.IsDir() {
					return nil
				}

				return errors.Errorf("%s already exists as file but it should be directory", outPath)
			}

			err = os.Mkdir(outPath, 0o777)

			return errors.Wrapf(err, "failed to create folder %s", outPath)
		}

		tmpl, err := template.New(path.Base(entryPath)).
			Funcs(generator.templateFuncs()).
			ParseFS(s.Structure(), entryPath)
		if err != nil {
			return err
		}

		file, err := createFile(outPath, 0o666)
		if err != nil {
			return err
		}
		defer file.Close()

		view := &TemplateView{
			ProjectID:   opts.ProjectID,
			Description: s.Description(),
			MetaView:    MetaView{},
		}

		err = tmpl.Execute(file, view)
		if err != nil {
			return err
		}

		// If file intentionally marked as executable set up correct permissions
		if view.MetaView.executable {
			err = file.Chmod(0o776)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (generator *generator) findScratch(name string) (scratch.Scratch, error) {
	for _, scr := range generator.scratches {
		if scr.Name() == name {
			return scr, nil
		}
	}
	return nil, errors.Errorf("scratch %s not found", name)
}

func (generator *generator) updateOutPathByScratchRenameSpec(outPath string, s scratch.Scratch) string {
	directory, base := path.Split(outPath)

	if newBase, exists := s.RenameSpec()[base]; exists {
		return path.Join(directory, newBase)
	}

	return outPath
}

func createFile(name string, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
}
