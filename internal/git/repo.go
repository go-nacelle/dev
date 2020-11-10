package git

import (
	"fmt"
	"os/exec"
	"regexp"
)

var remoteURLPattern = regexp.MustCompile(`git@github\.com:go-nacelle/(.+)\.git`)

// Repo returns the name of the repo from the current working directory.
// The resulting name does not include the `go-nacelle/` prefix.
func Repo() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}

	if matches := remoteURLPattern.FindStringSubmatch(string(out)); len(matches) > 0 {
		return matches[1], nil
	}

	return "", fmt.Errorf("unrecognized remote URL pattern")
}
