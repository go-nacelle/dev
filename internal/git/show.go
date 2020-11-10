package git

import (
	"fmt"
	"os/exec"
)

// Show returns the contents of the given file path at the given revision.
func Show(path string, revision string) ([]byte, error) {
	return exec.Command("git", "show", fmt.Sprintf("%s:%s", revision, path)).Output()
}
