package cli

import (
	"os"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
)

type interactiveFile struct {
	content     string
	description string
}

func newEditFile(content any, description string) (*interactiveFile, error) {
	c, err := toml.Marshal(content)
	if err != nil {
		return nil, err
	}

	return &interactiveFile{content: string(c), description: description}, nil
}

func (ef *interactiveFile) parse() error {
	tmp, err := os.CreateTemp("", "todo-*.toml")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(ef.description + "\n\n"); err != nil {
		return err
	}

	if _, err := tmp.WriteString(ef.content); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	c, err := os.ReadFile(tmp.Name())
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(c, ef.content); err != nil {
		return err
	}

	return nil
}

func (c *cli) StartInteractiveMode(content any, description string) error {
	edit, err := newEditFile(content, description)
	if err != nil {
		return err
	}

	if err := edit.parse(); err != nil {
		return err
	}

	return nil
}
