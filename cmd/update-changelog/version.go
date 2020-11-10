package main

import (
	"github.com/go-nacelle/dev/internal/git"
	"github.com/go-nacelle/dev/internal/modfile"
)

func readVersions() (map[string]map[string]string, error) {
	tags, err := git.Tags()
	if err != nil {
		return nil, err
	}
	tags = append(tags, "master")

	versions := map[string]map[string]string{}
	for _, tag := range tags {
		modfileContents, err := git.Show("go.mod", tag)
		if err != nil {
			return nil, err
		}

		tagVersions, err := modfile.Parse(modfileContents)
		if err != nil {
			return nil, err
		}

		versions[tag] = tagVersions
	}

	return versions, nil
}
