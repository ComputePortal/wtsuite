package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type PreProc struct {
  TokenData
}

type Extension struct {
  extension string
  behavior string
  PreProc
}

type Version struct {
  version int
  es string
  PreProc
}

func newPreProc(ctx context.Context) PreProc {
  return PreProc{newTokenData(ctx)}
}

func NewExtension(extension string, behavior string, ctx context.Context) *Extension {
  return &Extension{extension, behavior, newPreProc(ctx)}
}

func NewVersion(version int, es string, ctx context.Context) *Version {
  return &Version{version, es, newPreProc(ctx)}
}

func (t *Extension) Dump(indent string) string {
  var b strings.Builder
  
  b.WriteString(indent)

  b.WriteString("#extension ")
  b.WriteString(t.extension)
  b.WriteString(":")
  b.WriteString(t.behavior)

  return b.String()
}

// TODO: only write when no in library
func (t *Extension) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString("#extension ")
  b.WriteString(t.extension)
  b.WriteString(":")
  b.WriteString(t.behavior)
  
  return b.String()
}

func (t *Version) Dump(indent string) string {
  var b strings.Builder
  
  b.WriteString(indent)

  b.WriteString("#version ")
  b.WriteString(strconv.Itoa(t.version))
  b.WriteString(" ")
  b.WriteString(t.es)

  return b.String()
}

// written by tree/shader/ShaderBundle instead
func (t *Version) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  return ""
}

func (t *Version) CollectVersion(version *Word) (*Word, error) {
  thisVersionStr := strconv.Itoa(t.version) + " " + t.es
  if version != nil {
    if thisVersionStr != version.Value() {
      errCtx := t.Context()
      err := errCtx.NewError("Error: version mismatch")
      err.AppendContextString("Info: version declared here", version.Context())
      return nil, err
    }

    return version, nil
  } else {
    return NewWord(thisVersionStr, t.Context()), nil
  }
}

func (t *PreProc) ResolveStatementNames(scope Scope) error {
  return nil
}

func (t *PreProc) EvalStatement() error {
  return nil
}

func (t *PreProc) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *PreProc) UniqueStatementNames(ns Namespace) error {
  return nil
}
