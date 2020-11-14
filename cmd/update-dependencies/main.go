package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"

	"github.com/blang/semver"
	"github.com/go-nacelle/dev/internal/git"
	"github.com/go-nacelle/dev/internal/modfile"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainErr() error {
	headModfileContents, err := git.Show("go.mod", "head")
	if err != nil {
		return err
	}
	headDependencyVersions, err := modfile.Parse(headModfileContents)
	if err != nil {
		return err
	}

	for name := range headDependencyVersions {
		dependency := fmt.Sprintf("github.com/go-nacelle/%s", name)

		if err := exec.Command("go", "get", "-u", dependency).Run(); err != nil {
			return err
		}
	}

	currentModfileContents, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return err
	}
	currentDependencyVersions, err := modfile.Parse(currentModfileContents)
	if err != nil {
		return err
	}

	changes := []int64{0, 0, 0}

	var names []string
	for name := range headDependencyVersions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		headVersion, err := semver.ParseTolerant(headDependencyVersions[name])
		if err != nil {
			return err
		}

		currentVersion, err := semver.ParseTolerant(currentDependencyVersions[name])
		if err != nil {
			return err
		}

		if currentVersion.Major != headVersion.Major {
			fmt.Printf("%s: %s -> %s (major version update)\n", name, headVersion, currentVersion)
			changes[0]++
		} else if currentVersion.Minor != headVersion.Minor {
			fmt.Printf("%s: %s -> %s (minor version update)\n", name, headVersion, currentVersion)
			changes[1]++
		} else if currentVersion.Patch != headVersion.Patch {
			fmt.Printf("%s: %s -> %s (patch version update)\n", name, headVersion, currentVersion)
			changes[2]++
		}
	}

	currentVersionStr, err := git.Tag()
	if err != nil {
		currentVersionStr = "0.0.0"
	}

	currentVersion, err := semver.ParseTolerant(currentVersionStr)
	if err != nil {
		return err
	}
	newVersion := currentVersion

	if changes[0] != 0 {
		newVersion.Major++
		newVersion.Minor = 0
		newVersion.Patch = 0

		fmt.Printf("\n")
		fmt.Printf("Bump major: %s -> %s\n", currentVersion, newVersion)
		return nil
	}

	if changes[1] != 0 {
		newVersion.Minor++
		newVersion.Patch = 0

		fmt.Printf("\n")
		fmt.Printf("Bump minor: %s -> %s\n", currentVersion, newVersion)
		return nil
	}

	if changes[2] != 0 {
		newVersion.Patch++

		fmt.Printf("\n")
		fmt.Printf("Bump patch: %s -> %s\n", currentVersion, newVersion)
		return nil
	}

	fmt.Printf("\n")
	fmt.Printf("No change\n")
	return nil
}
