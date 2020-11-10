package main

import (
	"io/ioutil"
	"os"
)

func rewriteFile(path string, f func(contents []byte) ([]byte, error)) error {
	contents, err := ioutil.ReadFile("CHANGELOG.md")
	if err != nil {
		return err
	}

	contents, err = f(contents)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("CHANGELOG.md", contents, os.ModePerm); err != nil {
		return err
	}

	return nil
}
