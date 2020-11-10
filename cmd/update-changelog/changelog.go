package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-nacelle/dev/internal/changelog"
)

func readChangelogs(versions map[string]map[string]string) (map[string]*changelog.Changelog, error) {
	changelogs := map[string]*changelog.Changelog{}
	for _, tagVersions := range versions {
		for repo := range tagVersions {
			if _, ok := changelogs[repo]; !ok {
				changelog, err := readChangelog(repo)
				if err != nil {
					return nil, err
				}

				changelogs[repo] = changelog
			}
		}
	}

	return changelogs, nil
}

func readChangelog(repo string) (*changelog.Changelog, error) {
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/go-nacelle/%s/master/CHANGELOG.md", repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	changelog, err := changelog.ParseChangelog(strings.Split(string(contents), "\n"))
	if err != nil {
		return nil, err
	}

	return changelog, nil
}
