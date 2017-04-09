package upgrade

import (
	"github.com/dnote-io/cli/utils"
	"gopkg.in/yaml.v2"
)

func isYAML(b []byte) bool {
	var note utils.YAMLDnote

	err := yaml.Unmarshal(b, &note)
	return err == nil
}

func isDnoteUsingYAML() (bool, error) {
	b, err := utils.ReadNoteContent()
	if err != nil {
		return false, err
	}

	return isYAML(b), nil
}
