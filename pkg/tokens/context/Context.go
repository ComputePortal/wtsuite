package context

import (
  "io/ioutil"
  "strings"
)

type Source struct {
	source string
}

func NewSource(src string) *Source {
	return &Source{src}
}

func (s *Source) GetChar(i int) byte {
  return s.source[i]
}

func (s *Source) Len() int {
  return len(s.source)
}

type Context struct {
	ranges []struct{ start, stop int }
	source *Source
	path   string // where context is defined
}

func newContext(start, stop int, source *Source, path string) Context {
	return Context{
		[]struct{ start, stop int }{{start, stop}},
		source,
		path,
	}
}

// for preset globals
func NewDummyContext() Context {
	return newContext(0, 0, &Source{""}, "")
}

func NewContext(source *Source, path string) Context {
	if path == "" {
		panic("use dummycontext instead")
	}
	return newContext(0, source.Len(), source, path)
}

func (c *Context) NewContext(relStart, relStop int) Context {
	start := c.ranges[0].start
	return newContext(start+relStart, start+relStop, c.source, c.path)
}

func (c *Context) getRange(i int) (int, int) {
	if i == -1 {
		return c.ranges[0].start, c.ranges[len(c.ranges)-1].stop
	} else {
		return c.ranges[i].start, c.ranges[i].stop
	}
}

func (c *Context) IsConsecutive(other Context) bool {
	_, stop := c.getRange(-1)
	start, _ := other.getRange(-1)
	return stop == start
}

func (c *Context) appendRange(start, stop int) {
	c.ranges = append(c.ranges, struct{ start, stop int }{start, stop})
}

func (c *Context) slice(start, stop int) Context {
	result := Context{
		[]struct{ start, stop int }{},
		c.source,
		c.path,
	}

	if a, b := c.getRange(-1); stop < a || start > b {
		return result
	}

	for _, r := range c.ranges {
		a, b := r.start, r.stop

		as, bs := -1, -1

		if b > start {
			if a < start {
				as = start
			} else if a < stop {
				as = a
			} else {
				break
			}

			if b < stop {
				bs = b
			} else {
				bs = stop
			}
			result.appendRange(as, bs)
		}
	}

	return result
}

func (a *Context) Merge(b Context) Context {
	if a.path != b.path {
		panic("not same file")
	}

	c := Context{
		[]struct{ start, stop int }{},
		a.source,
		a.path,
	}

	ia, na := 0, len(a.ranges)
	ib, nb := 0, len(b.ranges)
	start, stop := -1, -1

	for ia < na || ib < nb || start != -1 {
		if start == -1 {
			// start a new range
			if (ia < na) && ((ib >= nb) || (a.ranges[ia].start < b.ranges[ib].start)) {
				start, stop = a.getRange(ia)
				ia++
			} else {
				start, stop = b.getRange(ib)
				ib++
			}
		} else {
			if (ia < na) && (a.ranges[ia].start <= stop) {
				// extend the range using a
				if next := a.ranges[ia].stop; next > stop {
					stop = next
				}
				ia++
			} else if (ib < nb) && b.ranges[ib].start <= stop {
				// extend the range using b
				if next := b.ranges[ib].stop; next > stop {
					stop = next
				}
				ib++
			} else {
				// append the range, and reset
				c.appendRange(start, stop)
				start, stop = -1, -1
			}
		}
	}

	return c
}

func MergeContexts(cs ...Context) Context {
	all := cs[0]
	for _, c := range cs[1:] {
		all = all.Merge(c)
	}

	return all
}

func (c *Context) Less(other *Context) bool {
  a := c.ranges[0].start
  b := other.ranges[0].start
  if a == b {
    a = c.ranges[0].stop
    b = other.ranges[0].stop

    return a < b
  } else {
    return a < b 
  }
}

func (c *Context) Same(other *Context) bool {
  a := c.ranges[0].start
  b := other.ranges[0].start
  if a == b {
    a = c.ranges[0].stop
    b = other.ranges[0].stop

    return a == b
  } else {
    return false
  }
}

func (c *Context) IsAtLineStart() bool {
  i := c.ranges[0].start

  if i <= 0 {
    return true
  } else {
    char := (c.source.source)[i-1]
    if char == '\n' || char == '\r' {
      return true
    } else {
      return false
    }
  }
}

func (c *Context) IsAtSourceStart() bool {
  i := c.ranges[0].start

  if i <= 0 {
    return true
  }

  return false
}

func (c *Context) IsSingleLine() bool {
	start := c.ranges[0].start
	stop := c.ranges[len(c.ranges)-1].stop

  for i := start; i < stop; i++ {
    char := (c.source.source)[i]

    if char == '\n' || char == '\r' {
      return false
    }
  }

  return true
}

// 0 if stop == start
func (c *Context) Distance(other *Context) int {
  if c.ranges[0].start < other.ranges[0].start {
    return other.ranges[0].start - c.ranges[len(other.ranges)-1].stop
  } else {
    return c.ranges[0].start - other.ranges[len(other.ranges)-1].stop
  }
}

func (c *Context) Len() int {
	return c.ranges[len(c.ranges)-1].stop - c.ranges[0].start
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) Content() string {
	start := c.ranges[0].start
	stop := c.ranges[len(c.ranges)-1].stop

  //if stop < start {
    //start, stop = stop, start
  //}

	return (c.source.source)[start:stop]
}

// replace sadly means reading entire file, and then replacing it
// the reading has already been done though, so just use that
// WARNING: this is irreversible and should only be used during refactorings
func (c *Context) SearchReplaceOrig(old, new string) error {
  var b strings.Builder

  prev := 0
  for _, r := range c.ranges {
    b.WriteString((c.source.source)[prev:r.start])

    b.WriteString(strings.Replace((c.source.source)[r.start:r.stop], old, new, -1))

    prev = r.stop
  }

  if prev < c.source.Len() {
    b.WriteString((c.source.source)[prev:c.source.Len()])
  }

  if err := ioutil.WriteFile(c.Path(), []byte(b.String()), 0); err != nil {
    return err
  }

  return nil
}

func MergeFill(a Context, b Context) Context {
  ctx := a.Merge(b)

  ctx = newContext(ctx.ranges[0].start, ctx.ranges[len(ctx.ranges)-1].stop, ctx.source, 
    ctx.path)

  return ctx
}

// 
func SimpleFill(a Context, b Context) Context {
  if a.path != b.path {
    panic("must be in same file")
  }

  start := a.ranges[0].start
  stop := b.ranges[len(b.ranges)-1].stop

  if start >= stop {
    start, stop = stop, start
  }

  ctx := newContext(start, stop, a.source, a.path)

  return ctx
}

func (c *Context) IncludeLeftSpace() Context {
  start := c.ranges[0].start
  stop := c.ranges[len(c.ranges)-1].stop

  for start > 0 && c.source.GetChar(start-1) == ' ' {
    start -= 1
  }

  return newContext(start, stop, c.source, c.path)
}

func (c *Context) IncludeRightSpace() Context {
  start := c.ranges[0].start
  stop := c.ranges[len(c.ranges)-1].stop

  if stop < c.source.Len() - 1 && c.source.GetChar(stop) == ' ' {
    stop += 1
  }

  return newContext(start, stop, c.source, c.path)
}
