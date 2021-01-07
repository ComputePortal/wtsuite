package parsers

import (
  "errors"
  "path/filepath"
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/files"
)

type CLIOption interface {
  Short() string // doesnt include initial dash
  Long() string // gnu extension --
  NArgs() int // 0 for flag, 1 for regular option, nice messages for each arg
  Info() string // a short message explaining the option
  Handle([]string) error // callback
}



type CLIOptionData struct {
  done bool
  unique bool
  short string
  long  string
  info  string
}

func newCLIOptionData(unique bool, short, long, info string) CLIOptionData {
  return CLIOptionData{false, unique, short, long, info}
}

func (o *CLIOptionData) Short() string {
  return o.short
}

func (o *CLIOptionData) Long() string {
  return o.long
}

func (o *CLIOptionData) Info() string {
  return o.info
}

func (o *CLIOptionData) Handle(args []string) error {
  if o.unique {
    if o.done {
      return errors.New("already done")
    }
  }

  o.done = true

  return nil
}



type CLIFlag struct {
  target *bool
  CLIOptionData
}

func NewCLIFlag(unique bool, short, long, info string, target *bool) CLIOption {
  return &CLIFlag{target, newCLIOptionData(false, short, long, info)}
}

func NewCLIUniqueFlag(short, long, info string, target *bool) CLIOption {
  return &CLIFlag{target, newCLIOptionData(true, short, long, info)}
}

func (o *CLIFlag) NArgs() int {
  return 0
}

func (o *CLIFlag) Handle(args []string) error {
  // could already have been set before
  *(o.target) = true

  return o.CLIOptionData.Handle(args)
}



type CLICountFlag struct {
  target *int
  CLIOptionData
}

func NewCLICountFlag(short, long, info string, target *int) CLIOption {
  return &CLICountFlag{target, newCLIOptionData(false, short, long, info)}
}

func (o *CLICountFlag) NArgs() int {
  return 0
}

func (o *CLICountFlag) Handle(args []string) error {
  *(o.target) = *(o.target) + 1

  return o.CLIOptionData.Handle(args)
}



type CLIString struct {
  target *string
  CLIOptionData
}

func newCLIString(unique bool, short, long, info string, target *string) CLIString {
  return CLIString{target, newCLIOptionData(unique, short, long, info)}
}

func NewCLIString(short, long, info string, target *string) CLIOption {
  o := newCLIString(false, short, long, info, target)
  return &o
}

func NewCLIUniqueString(short, long, info string, target *string) CLIOption {
  o := newCLIString(true, short, long, info, target)
  return &o
}

func (o *CLIString) NArgs() int {
  return 1
}

func (o *CLIString) Handle(args []string) error {
  *(o.target) = args[0]

  return o.CLIOptionData.Handle(args)
}



type CLIInt struct {
  target *int
  CLIOptionData
}

func NewCLIInt(short, long, info string, target *int) CLIOption {
  return &CLIInt{target, newCLIOptionData(false, short, long, info)}
}

func NewCLIUniqueInt(short, long, info string, target *int) CLIOption {
  return &CLIInt{target, newCLIOptionData(true, short, long, info)}
}

func (o *CLIInt) NArgs() int {
  return 1
}

func (o *CLIInt) Handle(args []string) error {
  i, err := strconv.ParseInt(args[0], 10, 64)
  if err != nil {
    return errors.New("invalid int")
  }

  *(o.target) = int(i)

  return o.CLIOptionData.Handle(args)
}



type CLIFile struct {
  checkExistence bool
  CLIString
}

func NewCLIFile(short, long, info string, checkExistence bool, target *string) CLIOption {
  return &CLIFile{checkExistence, newCLIString(false, short, long, info, target)}
}

func NewCLIUniqueFile(short, long, info string, checkExistence bool, target *string) CLIOption {
  return &CLIFile{checkExistence, newCLIString(true, short, long, info, target)}
}

func (o *CLIFile) Handle(args []string) error {
  path, err := filepath.Abs(args[0])
  if err != nil {
    panic(err)
    return err
  }

  // check that is an actual file, and not a directory
  if o.checkExistence {
    if err := files.AssertFile(path); err != nil {
      return err
    } 
  }

  return o.CLIString.Handle([]string{path})
}



type CLIDir struct {
  checkExistence bool
  CLIString
}

func NewCLIDir(short, long, info string, checkExistence bool, target *string) CLIOption {
  return &CLIDir{checkExistence, newCLIString(false, short, long, info, target)}
}

func NewCLIUniqueDir(short, long, info string, checkExistence bool, target *string) CLIOption {
  return &CLIDir{checkExistence, newCLIString(true, short, long, info, target)}
}

func (o *CLIDir) Handle(args []string) error {
  path, err := filepath.Abs(args[0])
  if err != nil {
    return err
  }

  // check that is an actual file, and not a directory
  if o.checkExistence {
    if err := files.AssertDir(path); err != nil {
      return err
    } 
  }

  return o.CLIString.Handle([]string{path})
}



type CLIEnum struct {
  valid []string
  CLIString
}

func NewCLIEnum(short, long, info string, valid []string, target *string) CLIOption {
  return &CLIEnum{valid, newCLIString(false, short, long, info, target)}
}

func NewCLIUniqueEnum(short, long, info string, valid []string, target *string) CLIOption {
  return &CLIEnum{valid, newCLIString(true, short, long, info, target)}
}

func (o *CLIEnum) Handle(args []string) error {
  found := ""

  for _, v := range o.valid {
    if strings.HasPrefix(v, args[0]) {
      if found != "" {
        return errors.New("ambiguous (\"" + args[0] + "\" is \"" + found + "\" or \"" + v + "\"?)")
      }

      found = v
    }
  }

  return o.CLIString.Handle([]string{found})
}



type CLICustom struct {
  nargs int
  handle func([]string) error
  CLIOptionData
}

func NewCLICustom(short, long, info string, nargs int, handle func([]string) error) CLIOption {
  return &CLICustom{nargs, handle, newCLIOptionData(false, short, long, info)}
}

func NewCLIUniqueCustom(short, long, info string, nargs int, handle func([]string) error) CLIOption {
  return &CLICustom{nargs, handle, newCLIOptionData(true, short, long, info)}
}

func (o *CLICustom) NArgs() int {
  return o.nargs
}

func (o *CLICustom) Handle(args []string) error {
  return o.handle(args)
}



type CLIAppendString struct {
  target *[]string
  CLIOptionData
}

func NewCLIAppendString(short, long, info string, target *[]string) CLIOption {
  return &CLIAppendString{target, newCLIOptionData(false, short, long, info)}
}

func (o *CLIAppendString) NArgs() int {
  return 1
}

func (o *CLIAppendString) Handle(args []string) error {
  *(o.target) = append(*(o.target), args[0])

  return o.CLIOptionData.Handle(args)
}



type CLIKeyValue struct {
  target map[string]string
  CLIOptionData
}

func NewCLIKeyValue(short, info string, target map[string]string) CLIOption {
  return &CLIKeyValue{target, newCLIOptionData(false, short, "", info)}
}

func NewCLIUniqueKeyValue(short, info string, target map[string]string) CLIOption {
  return &CLIKeyValue{target, newCLIOptionData(true, short, "", info)}
}

func (o *CLIKeyValue) NArgs() int {
  return 2
}

func (o *CLIKeyValue) Handle(args []string) error {
  if prev, ok := o.target[args[0]]; ok && o.unique {
    return errors.New("key " + args[0] + " already set to " + prev)
  }

  o.target[args[0]] = args[1]

  o.done = true

  return nil
}



type CLIKey struct {
  target map[string]string // fills with empty strings
  CLIOptionData
}

func NewCLIKey(short, info string, target map[string]string) CLIOption {
  return &CLIKey{target, newCLIOptionData(false, short, "", info)}
}

func NewCLIUniqueKey(short, info string, target map[string]string) CLIOption {
  return &CLIKey{target, newCLIOptionData(true, short, "", info)}
}

func (o *CLIKey) NArgs() int {
  return 1
}

func (o *CLIKey) Handle(args []string) error {
  if _, ok := o.target[args[0]]; ok && o.unique {
    return errors.New("key " + args[0] + " already set")
  }

  o.target[args[0]] = ""

  o.done = true

  return nil
}



// use as last
type CLIRemaining struct {
  target *[]string
  CLIOptionData
}

func NewCLIRemaining(target *[]string) CLIOption {
  return &CLIRemaining{target, newCLIOptionData(false, "", "", "")}
}

func (o *CLIRemaining) NArgs() int {
  return -1
}

func (o *CLIRemaining) Handle(args []string) error {
  *(o.target) = args

  return o.CLIOptionData.Handle(args)
}
