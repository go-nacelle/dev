package git

import (
	"os/exec"
	"strings"
)

// Tags returns a list of git tags in the current working directory.
func Tags() ([]string, error) {
	out, err := exec.Command("git", "tag").Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}
