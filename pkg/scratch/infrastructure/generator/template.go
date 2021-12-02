package generator

import (
	"io"
	"strings"
)

const (
	projectIDPathVar = "{{ProjectID}}"
)

type TemplateView struct {
	ProjectID   string
	Description string

	MetaView MetaView
}

type MetaView struct {
	executable bool
}

func (v *MetaView) MarkAsExecutable() bool {
	v.executable = true
	return true
}

func (generator *generator) templateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"toLowerCamelCase": func(value string) string {
			if value == "" {
				return ""
			}
			firstElemInLower := strings.ToLower(string(value[0]))
			if len(value) == 1 {
				return firstElemInLower
			}
			return firstElemInLower + value[1:]
		},
		"toUpperCamelCase": func(value string) string {
			if value == "" {
				return ""
			}
			firstElemInUpper := strings.ToUpper(string(value[0]))
			if len(value) == 1 {
				return firstElemInUpper
			}
			return firstElemInUpper + value[1:]
		},
		"toGOString": toGoString,
		"import": func(scratchName string, filePath string) (string, error) {
			s, err := generator.findScratch(scratchName)
			if err != nil {
				return "", err
			}

			file, err := s.Structure().Open(filePath)
			if err != nil {
				return "", err
			}

			bytes, err := io.ReadAll(file)
			if err != nil {
				return "", err
			}

			return string(bytes), nil
		},
	}
}

func toGoString(value string) string {
	return strings.NewReplacer("-", "", "_", "").Replace(value)
}

func replaceTemplateVariablesInPath(outPath, projectID string) string {
	// Since there is only one variable ProjectID just replace it
	return strings.ReplaceAll(outPath, projectIDPathVar, projectID)
}
