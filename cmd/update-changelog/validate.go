package main

import (
	"fmt"

	"github.com/go-nacelle/dev/internal/changelog"
)

func validateTags(changelog *changelog.Changelog, versions map[string]map[string]string) error {
	var versionStrings []string
	for _, version := range changelog.Versions {
		if _, ok := versions[version.Version]; !ok && version.Version != "Unreleased" {
			return fmt.Errorf("version %s is not a git tag", version.Version)
		}

		versionStrings = append(versionStrings, version.Version)
	}

	for version := range versions {
		var found bool
		for _, x := range versionStrings {
			if x == version || (x == "Unreleased" && version == "master") {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("git tag %s is not a version", version)
		}
	}

	return nil
}
