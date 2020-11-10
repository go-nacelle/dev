package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseChangelog converts an existing changelog's contnets split by newline
// into a parsed changelog object. The prefix (content before the first version)
// and suffix (content after ---) is also preserved.
func ParseChangelog(lines []string) (*Changelog, error) {
	versions := parseVersions(lines)
	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found in changelog")
	}

	prefix := strings.TrimSpace(strings.Join(lines[:versions[0].LineOffset], "\n"))
	if prefix != "" {
		prefix += "\n\n"
	}

	suffix := ""
	n := versions[len(versions)-1].LineOffset
	for i, line := range lines[n:] {
		if line == "---" {
			// Anything after --- is considered the suffix
			suffix = strings.Join(lines[n+i+1:], "\n")
			break
		}
	}

	changelog := &Changelog{
		Prefix:   prefix,
		Suffix:   suffix,
		Versions: versions,
	}

	return changelog, nil
}

type parsedVersion struct {
	offset     int
	version    string
	releasedOn string
}

var h2Pattern = regexp.MustCompile(`^## \[(.+)\](?: - (\d{4}-\d{2}-\d{2}))?$`)

// parseVersions converts the given file contents into an order-preserving
// list of parsed release versions.
func parseVersions(lines []string) []*Version {
	var parsedVersions []parsedVersion
	for i, line := range lines {
		if match := h2Pattern.FindStringSubmatch(line); len(match) > 0 {
			parsedVersions = append(parsedVersions, parsedVersion{
				offset:     i,
				version:    string(match[1]),
				releasedOn: string(match[2]),
			})
		}
	}

	var versions []*Version
	for i, version := range parsedVersions {
		nextOffset := len(lines) - 1
		if i != len(parsedVersions)-1 {
			nextOffset = parsedVersions[i+1].offset
		}

		versions = append(versions, &Version{
			Version:      version.version,
			ReleasedOn:   version.releasedOn,
			ChangeGroups: parseChangeGroups(lines[version.offset+1 : nextOffset]),
			LineOffset:   version.offset,
		})
	}

	return versions
}

type parsedChangeGroup struct {
	offset     int
	changeType string
}

var h3Pattern = regexp.MustCompile(`^### (.+)$`)
var downstreamChangePattern = regexp.MustCompile(`^\[go-nacelle/.+@.+\] -> \[go-nacelle/.+@.+\]$`)

// parseChangeGroups converts the given file content between h2 elements in
// the changelog text into an ordered-preserving list of change groups for a
// particular version.
func parseChangeGroups(lines []string) (changeGroups []*ChangeGroup) {
	var groups []parsedChangeGroup
	for i, line := range lines {
		if match := h3Pattern.FindStringSubmatch(line); len(match) > 0 {
			groups = append(groups, parsedChangeGroup{
				offset:     i,
				changeType: string(match[1]),
			})
		}
	}

	for i, group := range groups {
		nextOffset := len(lines)
		if i != len(groups)-1 {
			nextOffset = groups[i+1].offset
		}

		var changes []Renderable
		for _, line := range lines[group.offset+1 : nextOffset] {
			// If we hit ---, stop parsing the remainder of this version.
			// This is an indicator that there is footer information in
			// the changelog that should not be parsed as automatically
			// formatted changelog content.
			if line == "---" {
				return changeGroups
			}

			if strings.HasPrefix(line, "- ") {
				line := strings.TrimSpace(line[1:])

				if !downstreamChangePattern.MatchString(line) {
					changes = append(changes, &Change{
						Description: line,
					})
				}
			}
		}

		changeGroups = append(changeGroups, &ChangeGroup{
			ChangeType: group.changeType,
			Changes:    changes,
		})

	}

	return changeGroups
}
