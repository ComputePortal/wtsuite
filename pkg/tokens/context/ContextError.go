package context

import (
	"os"
	"path/filepath"
	"strings"
)

type ContextError struct {
	err string
}

// exported because also used in files modules
func Abbreviate(path string) string {
	pwd, err := os.Getwd()
	if err != nil {
		return path
	}

	pwd, err = filepath.Abs(pwd)
	if err != nil {
		return path
	}

	relPath, err := filepath.Rel(pwd, path)
	if err != nil {
		return path
	}

	return relPath
}

func (ce *ContextError) Error() string {
	return ce.err
}

func (c *Context) NewError(msg string) *ContextError {
	ce := &ContextError{""}
	ce.AppendContextString(msg, *c)

	return ce
}

func (ce *ContextError) prepareContextString(msg string, c Context) string {
	var b strings.Builder

	b.WriteString("\u001b[1m")
	b.WriteString(msg)
	b.WriteString(" \u001b[0m(\u001b[35m")
	b.WriteString(Abbreviate(c.path))
	b.WriteString("\u001b[0m)\n")

	cl := c.newContextLines()
	cl.pad(1)

	b.WriteString(cl.write(cl.lineNumberFormat("  ")))

	return b.String()
}

func (ce *ContextError) AppendContextString(msg string, c Context) {
	s := ce.prepareContextString(msg, c)

	if ce.err != "" {
		s = "\n" + s
	}

	if !strings.HasSuffix(ce.err, s) {
		ce.err += s
	}
}

func (ce *ContextError) AppendError(other *ContextError) {
  ce.err += "\n" + other.err
}

func (ce *ContextError) PrependContextString(msg string, c Context) {
	s := ce.prepareContextString(msg, c)

	if ce.err != "" {
		s += "\n"
	}

	if !strings.HasPrefix(ce.err, s) {
		ce.err = s + ce.err
	}
}

func (ce *ContextError) AppendString(msg string) {
	var b strings.Builder

	//b.WriteString("\u001b[1m")
	b.WriteString(msg)
	//b.WriteString("\u001b[0m\n")

	if !strings.HasSuffix(ce.err, b.String()) {
		ce.err += b.String()
	}
}

func (ce *ContextError) ToHTML() string {
	// create a valid html document
	r := strings.NewReader(ce.err)
	var w strings.Builder
	w.WriteString("<!doctype html><html><head><style>i{color:#f00; font-weight:bold}</style></head><body>")

	escaping := false
	escapeCode := ""
	activeTag := ""

	// simple bold uses b

	for true {
		if ch, _, err := r.ReadRune(); err == nil {
			if !escaping {
				switch string(ch) {
				case "\u001b":
					escaping = true
				case "<":
					w.WriteString("&lt;")
				case ">":
					w.WriteString("&gt;")
				case "\n":
					w.WriteString("<br>")
				default:
					w.WriteString(string(ch))
				}
			} else {
				if string(ch) == "m" {
					escaping = false

					if activeTag != "" {
						if escapeCode != "[0" {
							panic("unmatched escape")
						}

						w.WriteString("</")
						w.WriteString(activeTag)
						w.WriteString(">")
						activeTag = ""
					} else {
						switch escapeCode {
						case "[31;1":
							activeTag = "i"
						case "[1":
							activeTag = "b"
						default:
							panic("unexpected escapeCode " + escapeCode)
						}
						w.WriteString("<")
						w.WriteString(activeTag)
						w.WriteString(">")
					}

					escapeCode = ""
				} else {
					escapeCode += string(ch)
				}
			}
		} else {
			break
		}
	}

	w.WriteString("</body></html>")
	return w.String()
}

func AppendContextString(err error, msg string, c Context) {
	switch e := err.(type) {
	case *ContextError:
		e.AppendContextString(msg, c)
	}
}

func AppendError(err_ error, other_ error) {
  err, ok := err_.(*ContextError)
  if ok {
    other, ok := other_.(*ContextError)
    if ok {
      err.AppendError(other)
    }
  }
}

func PrependContextString(err error, msg string, c Context) {
	switch e := err.(type) {
	case *ContextError:
		e.PrependContextString(msg, c)
	}
}

func AppendString(err error, msg string) {
	switch e := err.(type) {
	case *ContextError:
		e.AppendString(msg)
	}
}

func ToHTML(err error) string {
	switch e := err.(type) {
	case *ContextError:
		return e.ToHTML()
	default:
		return e.Error()
	}
}
