package scratch

import (
	"bytes"
	"io"
	"io/fs"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Loader interface {
	Load(manifestPath string) ([]Scratch, error)
}

func NewFsLoader(f fs.FS) Loader {
	return &fsLoader{fs: f}
}

type fsLoader struct {
	fs fs.FS
}

func (loader *fsLoader) Load(manifestPath string) ([]Scratch, error) {
	manifestFile, err := loader.fs.Open(manifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open manifest file %s", manifestPath)
	}
	defer manifestFile.Close()

	manifestBytes, err := io.ReadAll(manifestFile)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load manifest file %s", manifestPath)
	}

	manifestsMap, err := loader.loadScratchesManifests(manifestBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse manifest %s", manifestPath)
	}

	file, err := loader.fs.Open(".")
	if err != nil {
		return nil, err
	}

	dir, ok := file.(fs.ReadDirFile)
	if !ok {
		return nil, errors.New("mapped wrong path into mutablefs.FS")
	}

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	if err = validateEntries(entries, manifestsMap); err != nil {
		return nil, err
	}

	scratches := make([]Scratch, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifest, exists := manifestsMap[entry.Name()]
		if !exists {
			continue
		}

		scratchFS, err2 := fs.Sub(loader.fs, manifest.Name)
		if err2 != nil {
			return nil, err2
		}

		scratches = append(scratches, &scratch{
			name:        manifest.Name,
			description: manifest.Description,
			fs:          scratchFS,
			renameSpec:  manifest.RenameSpec,
		})
	}

	return scratches, nil
}

type scratchManifest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	RenameSpec  map[string]string `yaml:"rename,omitempty"`
}

func (loader *fsLoader) loadScratchesManifests(resources []byte) (map[string]scratchManifest, error) {
	dec := yaml.NewDecoder(bytes.NewReader(resources))

	res := map[string]scratchManifest{}
	for {
		var value scratchManifest
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res[value.Name] = value
	}
	return res, nil
}

func validateEntries(dirEntries []fs.DirEntry, manifestsMap map[string]scratchManifest) error {
	for name := range manifestsMap {
		var found bool
		for _, entry := range dirEntries {
			if entry.IsDir() && entry.Name() == name {
				found = true
			}
		}
		if !found {
			return errors.Errorf("dir for %s not found", name)
		}
	}
	return nil
}
