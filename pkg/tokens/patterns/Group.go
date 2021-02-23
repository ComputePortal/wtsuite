package patterns

import (
	"regexp"
)

type Group interface {
	Start() string
	Stop() string
	MatchStart(s string) bool
	MatchStop(s string) bool
	StartStopRegexp() *regexp.Regexp // XXX: don't allow access to this
	IsTagGroup() bool
}

type groupData struct {
	start           string
	stop            string
	startRegexp     *regexp.Regexp
	stopRegexp      *regexp.Regexp
	startStopRegexp *regexp.Regexp
}

type tagGroup struct {
	startRegexp     *regexp.Regexp
	stopRegexp      *regexp.Regexp
	startStopRegexp *regexp.Regexp
}

type scriptTagGroup struct {
  inSingleQuote bool
  inDoubleQuote bool
  tagGroup
}

func NewGroup(start, stop string) *groupData {
	if start == "" {
		panic("start cant be empty") // except for NewTagGroup constructor
	}
	if stop == "" {
		panic("stop cant be empty") // except for NewTagGroup constructor
	}
	return &groupData{start, stop, compileRegexp(start), compileRegexp(stop), compileRegexp(start, stop)}
}

func newTagGroup(name string) tagGroup {
	//
	start := regexp.MustCompile(`<[\s]*` + name + `([\s][^</]*)?>`)
	stop := regexp.MustCompile(`<[\s]*[/]` + name + `[\s]*>`)
	startStop := regexp.MustCompile(`(` + start.String() + `)|(` + stop.String() + `)`)

	return tagGroup{start, stop, startStop}
}

func NewTagGroup(name string) *tagGroup {
  tg := newTagGroup(name)
  return &tg
}

func NewScriptTagGroup(name string) *scriptTagGroup {
  tg := newTagGroup(name)

  tg.startStopRegexp = regexp.MustCompile(tg.startStopRegexp.String() + `|(['"])`)

  return &scriptTagGroup{false, false, tg}
}

func (g *groupData) Start() string {
	if g.start == "" {
		panic("not available as key")
	}

	return g.start
}

func (g *tagGroup) Start() string {
	panic("not available as key")
}

func (g *groupData) Stop() string {
	if g.stop == "" {
		panic("not available as key")
	}

	return g.stop
}

func (g *tagGroup) Stop() string {
	panic("not available as key")
}

func (g *groupData) MatchStart(s string) bool {
	return g.startRegexp.MatchString(s)
}

func (g *groupData) MatchStop(s string) bool {
	return g.stopRegexp.MatchString(s)
}

func (g *tagGroup) MatchStart(s string) bool {
	return g.startRegexp.MatchString(s)
}

func (g *tagGroup) MatchStop(s string) bool {
	return g.stopRegexp.MatchString(s)
}

func (g *scriptTagGroup) matchInternal(s string, re *regexp.Regexp) bool {
  if g.inSingleQuote {
    if s == "'"  {
      g.inSingleQuote = false
    }

    return false
  } else if g.inDoubleQuote {
    if s =="\"" {
      g.inDoubleQuote = false
    }

    return false
  } else if s == "'" {
    g.inSingleQuote = true
    return false
  } else if s == "\"" {
    g.inDoubleQuote = true
    return false
  } else {
    return re.MatchString(s)
  }
}

func (g *scriptTagGroup) MatchStart(s string) bool {
  return g.matchInternal(s, g.startRegexp)
}

func (g *scriptTagGroup) MatchStop(s string) bool {
  return g.matchInternal(s, g.stopRegexp)
}

func (g *groupData) StartStopRegexp() *regexp.Regexp {
	return g.startStopRegexp
}

func (g *tagGroup) StartStopRegexp() *regexp.Regexp {
	return g.startStopRegexp
}

func (g *groupData) IsTagGroup() bool {
	return false
}

func (g *tagGroup) IsTagGroup() bool {
	return true
}
