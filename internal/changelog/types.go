package changelog

type Changelog struct {
	Prefix   string     // content that occurs before the latest version
	Suffix   string     // content that occurs after the first version
	Versions []*Version // reverse-chronological description of versions
}

type Version struct {
	Version      string       // version name
	ReleasedOn   string       // date version was released
	ChangeGroups ChangeGroups // grouped summary of changes
	LineOffset   int          // offset of header in original file
}

func (v *Version) AddChangeGroup(changeGroup *ChangeGroup) {
	v.ChangeGroups = v.ChangeGroups.AddChangeGroup(changeGroup)
}

type ChangeGroup struct {
	ChangeType string       // type of change (e.g. added, removed, changed)
	Changes    []Renderable // list of changes
}

type ChangeGroups []*ChangeGroup

func (changeGroups ChangeGroups) AddChangeGroup(changeGroup *ChangeGroup) ChangeGroups {
	for _, cg := range changeGroups {
		if cg.ChangeType == changeGroup.ChangeType {
			cg.Changes = append(cg.Changes, changeGroup.Changes...)
			return changeGroups
		}
	}

	return append(changeGroups, changeGroup)
}

type DependencyChange struct {
	DependencyName string       // nacelle dependency name
	OldVersion     string       // old dependency version
	NewVersion     string       // current dependency version
	Changes        []Renderable // list of changes
}

type Change struct {
	Description string // free-form text describing change
}
