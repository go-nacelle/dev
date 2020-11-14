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

// Tag returns the most recent git tag in the current working directory.
func Tag() (string, error) {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
