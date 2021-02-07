package styles

import (
  "errors"
  "io/ioutil"
  "strings"

	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func BuildDict(d *tokens.StringDict) (Sheet, error) {
  sheet := NewSheet()

  // the top level rules are added to the sheet
  if err := d.Loop(func(key *tokens.String, value_ tokens.Token, last bool) error {
    value, err := tokens.AssertStringDict(value_)
    if err != nil {
      return err
    }

    if strings.HasPrefix(key.Value(), "@") {
      atRules, err := ExpandAtRules(nil, key, value)
      if err != nil {
        return err
      }

      for _, r := range atRules {
        sheet.Append(r)
      }
    } else {
      sels, err := ParseSelectorList(key)
      if err != nil {
        return err
      }

      for _, sel := range sels{
        r := NewRule(sel, value)

        sheet.Append(r)
      }
    }

    return nil
  }); err != nil {
    return nil, err
  }

  expandedSheet, err := sheet.ExpandNested()
  if err != nil {
    return nil, err
  }

  return expandedSheet, nil
}

// expects export var style = {...} somewhere in file 
func Build(path string, ctx context.Context) (Sheet, error) {
  // always rebuild the 
  cache := directives.NewFileCache()

  scope, _, err := directives.BuildFile(cache, path, false, nil)
  if err != nil {
    return nil, err
  }

  if !scope.HasVar("main") {
    errCtx := ctx
    return nil, errCtx.NewError("Error: style var \"main\" not found in \"" + path + "\"")
  }

  v := scope.GetVar("main")
  if !v.Exported {
    errCtx := ctx
    return nil, errCtx.NewError("Error: style var \"main\" not exported from \"" + path + "\"")
  }

  // XXX: can functions with all defaults also be used?

  d, err := tokens.AssertStringDict(v.Value)
  if err != nil {
    return nil, err
  }

  return BuildDict(d)
}

func BuildFile(input string, output string) error {
  sheet, err := Build(input, context.NewDummyContext())
  if err != nil {
    return err
  }

  content, err := sheet.Write(true, patterns.NL, patterns.TAB)
  if err != nil {
    return err
  }

  if err := ioutil.WriteFile(output, []byte(content), 0644); err != nil {
    return errors.New("Error: " + err.Error())
  }

  return nil
}

func BuildDictWriteSheet(d *tokens.StringDict) (string, error) {
  sheet, err := BuildDict(d)
  if err != nil {
    return "", err
  }

  return sheet.Write(true, patterns.NL, patterns.TAB)
}

var _buildDictWriteSheetRegistered = directives.RegisterBuildStyle(BuildDictWriteSheet)
