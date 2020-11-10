package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-nacelle/dev/internal/changelog"
	"github.com/go-nacelle/dev/internal/git"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainErr() error {
	repo, err := git.Repo()
	if err != nil {
		return err
	}

	versions, err := readVersions()
	if err != nil {
		return err
	}

	changelogs, err := readChangelogs(versions)
	if err != nil {
		return err
	}

	return rewriteFile("CHANGELOG.md", func(contents []byte) ([]byte, error) {
		changelog, err := changelog.ParseChangelog(strings.Split(string(contents), "\n"))
		if err != nil {
			return nil, err
		}

		if err := validateTags(changelog, versions); err != nil {
			return nil, err
		}

		upstreamDependencyChangelogs2(changelog, changelogs, versions)
		return []byte(changelog.Render(repo)), nil
	})
}
