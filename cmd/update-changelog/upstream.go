package main

import (
	"sort"

	"github.com/go-nacelle/dev/internal/changelog"
)

func upstreamDependencyChangelogs2(candidateChangelog *changelog.Changelog, changelogs map[string]*changelog.Changelog, versions map[string]map[string]string) {
	for i, version := range candidateChangelog.Versions {
		if i+1 >= len(candidateChangelog.Versions) {
			// Skip last iteration
			break
		}

		v1 := version.Version
		if version.Version == "Unreleased" {
			v1 = "master"
		}
		v2 := candidateChangelog.Versions[i+1].Version

		var repos []string
		for repo := range changelogs {
			repos = append(repos, repo)
		}
		sort.Strings(repos)

		for _, repo := range repos {
			prev := versions[v2][repo]
			next := versions[v1][repo]
			if next != prev {
				upstreamDependencyChangelogs(version, changelogs, repo, prev, next)
			}
		}
	}
}

func upstreamDependencyChangelogs(version *changelog.Version, changelogs map[string]*changelog.Changelog, repo, prev, next string) {
	var versionedGroups changelog.ChangeGroups
	for _, version2 := range extractVersionsBetween(changelogs[repo].Versions, prev, next) {
		for _, group := range version2.ChangeGroups {
			versionedGroups = versionedGroups.AddChangeGroup(group)
		}
	}

	for _, cg := range versionedGroups {
		version.AddChangeGroup(&changelog.ChangeGroup{
			ChangeType: cg.ChangeType,
			Changes: []changelog.Renderable{
				&changelog.DependencyChange{
					DependencyName: repo,
					OldVersion:     prev,
					NewVersion:     next,
					Changes:        cg.Changes,
				},
			},
		})
	}
}

func extractVersionsBetween(versions []*changelog.Version, prev, next string) []*changelog.Version {
	pre := true
	var versionsBetween []*changelog.Version
	for _, v := range versions {
		if v.Version == next {
			pre = false
		}
		if pre {
			continue
		}
		if v.Version == prev {
			break
		}

		versionsBetween = append(versionsBetween, v)
	}

	return versionsBetween
}
