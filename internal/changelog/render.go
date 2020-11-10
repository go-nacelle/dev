package changelog

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func (c *Changelog) Render(repository string) string {
	var buf = bytes.Buffer{}
	c.renderInto(&buf, repository)
	return buf.String()
}

func (c *Changelog) renderInto(buf *bytes.Buffer, repository string) {
	context := &RenderContext{
		Links: map[string]string{},
	}

	for i := 0; i < len(c.Versions)-1; i++ {
		vx := c.Versions[i].Version
		if i == 0 {
			vx = "HEAD"
		}
		context.Links[c.Versions[i].Version] = fmt.Sprintf("https://github.com/go-nacelle/%s/compare/%s...%s", repository, c.Versions[i+1].Version, vx)
	}
	n := len(c.Versions) - 1
	context.Links[c.Versions[n].Version] = fmt.Sprintf("https://github.com/go-nacelle/%s/releases/tag/%s", repository, c.Versions[n].Version)

	fmt.Fprintf(buf, c.Prefix)

	for _, version := range c.Versions {
		version.renderInto(buf, context)
	}

	var keys []string
	for key := range context.Links {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintf(buf, "[%s]: %s\n", key, context.Links[key])
	}

	if c.Suffix != "" {
		fmt.Fprintf(buf, "\n---\n%s", c.Suffix)
	}
}

func (v *Version) renderInto(buf *bytes.Buffer, context *RenderContext) {
	if v.ReleasedOn == "" {
		fmt.Fprintf(buf, "## [%s]\n\n", v.Version)
	} else if v.ReleasedOn != "" {
		fmt.Fprintf(buf, "## [%s] - %s\n\n", v.Version, v.ReleasedOn)
	}

	for _, changeGroup := range v.ChangeGroups {
		changeGroup.renderInto(buf, context)
	}
}

func (g *ChangeGroup) renderInto(buf *bytes.Buffer, context *RenderContext) {
	empty := true
	for _, c := range g.Changes {
		if !c.isEmpty() {
			empty = false
			break
		}
	}
	if empty {
		return
	}

	fmt.Fprintf(buf, "### %s\n\n", g.ChangeType)
	for _, change := range g.Changes {
		change.renderInto(buf, 0, context)
	}
	fmt.Fprintf(buf, "\n")
}

type Renderable interface {
	isEmpty() bool
	renderInto(buf *bytes.Buffer, level int, context *RenderContext)
}

type RenderContext struct {
	Links map[string]string
}

func (c *Change) isEmpty() bool {
	return false
}

func (c *Change) renderInto(buf *bytes.Buffer, level int, context *RenderContext) {
	fmt.Fprintf(buf, "%s- %s\n", strings.Repeat(" ", level*2), c.Description)
}

func (c *DependencyChange) isEmpty() bool {
	for _, c := range c.Changes {
		if !c.isEmpty() {
			return false
		}
	}

	return true
}

func (c *DependencyChange) renderInto(buf *bytes.Buffer, level int, context *RenderContext) {
	if len(c.Changes) == 0 {
		return
	}

	n1 := fmt.Sprintf("go-nacelle/%s@%s", c.DependencyName, c.OldVersion)
	context.Links[n1] = fmt.Sprintf("https://github.com/go-nacelle/%s/releases/tag/%s", c.DependencyName, c.OldVersion)
	n2 := fmt.Sprintf("go-nacelle/%s@%s", c.DependencyName, c.NewVersion)
	context.Links[n2] = fmt.Sprintf("https://github.com/go-nacelle/%s/releases/tag/%s", c.DependencyName, c.NewVersion)

	fmt.Fprintf(buf, "%s- [%s] -> [%s]\n", strings.Repeat(" ", level*2), n1, n2)

	for _, change := range c.Changes {
		change.renderInto(buf, level+1, context)
	}
}
