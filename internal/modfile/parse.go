package modfile

import (
	"strings"

	"golang.org/x/mod/modfile"
)

const prefix = "github.com/go-nacelle/"

// Parse converts the given contents of a go.mod file into a mapping from
// go-nacelle repositories to that dependency's pinned version. The keys of
// the map do not include the `go-nacelle/` prefix.
func Parse(contents []byte) (map[string]string, error) {
	file, err := modfile.Parse("go.mod", contents, nil)
	if err != nil {
		return nil, err
	}

	versions := map[string]string{}
	for _, require := range file.Require {
		if strings.HasPrefix(require.Mod.Path, prefix) {
			versions[strings.TrimPrefix(require.Mod.Path, prefix)] = require.Mod.Version
		}
	}

	return versions, nil
}
