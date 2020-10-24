package context

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	HIGHLIGHT_START = "\u001b[31;1m" // bold red
	HIGHLIGHT_STOP  = "\u001b[0m"
)

type contextLines struct {
	ctx         *Context
	lines       []string
	active      []bool
	first, last int
}

func (c *Context) newContextLines() contextLines {
	lines := strings.Split(strings.Trim(c.source.source, "\n"), "\n")
	active := make([]bool, len(lines))

	first, last := -1, -1
	al, bl := 0, 0
	for il, line := range lines {
		bl = al + len(line) + 1 // also count the newline

		active[il] = len(c.slice(al, bl).ranges) > 0

		if first == -1 {
			first = il
		}
		last = il

		al = bl
	}

	return contextLines{
		c,
		lines,
		active,
		first, last,
	}
}

func (cl *contextLines) pad(np int) {
	nl := len(cl.lines)
	tmp := make([]bool, nl)

	for iter := 0; iter < np; iter++ {
		for i, a := range cl.active {
			if a {
				tmp[i] = true
			} else if i > 0 && cl.active[i-1] {
				tmp[i] = true
			} else if i < nl-1 && cl.active[i+1] {
				tmp[i] = true
			} else {
				tmp[i] = false
			}

			if tmp[i] {
				if i < cl.first {
					cl.first = i
				}
				if i > cl.last {
					cl.last = i
				}
			}
		}

		cl.active, tmp = tmp, cl.active
	}
}

func (cl *contextLines) loopLines(fn func(il, al, bl int, line string, active bool)) {
	al, bl := 0, 0
	for il, line := range cl.lines {
		bl = al + len(line)

		fn(il, al, bl, line, cl.active[il])

		al = bl + 1 // plus newline char (assuming no carriage returns)
	}
}

func (cl *contextLines) lineNumberFormat(prefix string) string {
	nd := math.Floor(math.Log10(float64(cl.last+1))) + 1

	return prefix + "\u001b[1m%0" + strconv.FormatInt(int64(nd), 10) + "d\u001b[0m "
}

func (cl *contextLines) write(lnf string) string {
	var b strings.Builder

	prevLine := -1

	cl.loopLines(func(il, al, bl int, line string, active bool) {
		c := cl.ctx.slice(al, bl)

		if active {
			if prevLine != -1 {
				if prevLine != il-1 {
					b.WriteString("\n  ...")
				}
				b.WriteString("\n")
			}

			b.WriteString(fmt.Sprintf(lnf, il+1)) // line indexing is 1-based

			prevStop := 0
			for _, r := range c.ranges {
				start, stop := r.start-al, r.stop-al
				b.WriteString(line[prevStop:start])
				b.WriteString(HIGHLIGHT_START)
				b.WriteString(line[start:stop])
				b.WriteString(HIGHLIGHT_STOP)
				prevStop = stop
			}

			if prevStop < len(line) {
				b.WriteString(line[prevStop:])
			}

			prevLine = il
		}
	})

	b.WriteString("\n")
	return b.String()
}

func (c *Context) WritePrettyOneLiner() string {
  cl := c.newContextLines()

  prefix := "\u001b[35m" + Abbreviate(c.path) + "\u001b[0m:"
  return cl.write(cl.lineNumberFormat(prefix))
}
