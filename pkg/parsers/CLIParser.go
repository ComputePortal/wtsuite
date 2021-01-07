package parsers

import (
  "errors"
  "runtime"
  "strconv"
  "strings"
)

const (
  SHORT_OPT_PREFIX = "-"
  LONG_OPT_PREFIX = "--"
  WIN_OPT_PREFIX = "/"
)

// very different from other (language) parsers
type CLIParser struct {
  args []string
  pos int
  info string
  postInfo string
  options []CLIOption
  last CLIOption // handles the remaining positional args
}

func NewCLIParser(info string, postInfo string, options []CLIOption, last CLIOption) *CLIParser {
  if last.NArgs() == 0 {
    panic("last option can't be a flag")
  }

  for _, opt := range options {
    if strings.HasPrefix(opt.Long(), SHORT_OPT_PREFIX) || strings.HasPrefix(opt.Long(), WIN_OPT_PREFIX) {
      panic("invalid long opt name \"" + opt.Long() + "\"")
    } else if strings.HasPrefix(opt.Short(), SHORT_OPT_PREFIX) || strings.HasPrefix(opt.Short(), WIN_OPT_PREFIX) {
      panic("invalid short opt name \"" + opt.Short() + "\"")
    }
  }

  return &CLIParser{nil, 0, info, postInfo, options, last}
}

func (p *CLIParser) Info() string {
  var b strings.Builder

  b.WriteString(p.info)
  b.WriteString("\n")

  b.WriteString("Options:\n")
  b.WriteString("  -?, -h, --help   Show this message\n")

  for _, option := range p.options {
    optInfo := option.Info()
    if optInfo == "" {
      continue
    }

    b.WriteString("  ")
    b.WriteString(optInfo)
  }

  if p.postInfo != "" {
    b.WriteString(p.postInfo)
    b.WriteString("\n")
  }

  return b.String()
}

// dont include the command name! (so os.Args()[1:])
func (p *CLIParser) start(args []string) {
  p.args = args
  p.pos = 0
}

func (p *CLIParser) End() bool {
  if p.args == nil {
    panic("not yet started")
  }

  return p.pos == len(p.args)
}

func (p *CLIParser) Eat(n int) ([]string, error) {
  if n == 0 {
    return []string{}, nil
  }

  if p.pos + n > len(p.args) {
    if p.pos == 0 {
      return nil, p.Error("expected at least " + strconv.Itoa(n) + " args")
    } else {
      return nil, p.Error("expected " + strconv.Itoa(n) + " args after " + p.args[p.pos-1])
    }
  }

  res := p.args[p.pos:p.pos+n]

  p.pos += n

  return res, nil
}

func (p *CLIParser) Error(msg string) error {
  var b strings.Builder

  b.WriteString("Error: ")
  b.WriteString(msg)
  b.WriteString("\n")
  b.WriteString(p.Info())

  return errors.New(b.String())
}

// returns nil if not found, or if ambiguous
func (p *CLIParser) findLongOption(key string) CLIOption {
  var res CLIOption = nil

  for _, opt := range p.options {
    if strings.HasPrefix(opt.Long(), key) {
      if res != nil {
        return nil
      }

      res = opt
    }
  }

  return res
}

func (p *CLIParser) findShortOption(key string) CLIOption {
  for _, opt := range p.options {
    if opt.Short() == key {
      return opt
    }
  }

  return nil
}

// returns the remaining positional arguments
func (p *CLIParser) Parse(args []string) error {
  p.start(args)

  positional := make([]string, 0)


  for !p.End() {
    arg_, err := p.Eat(1)
    if err != nil {
      return err
    }

    arg := arg_[0]

    if arg == "-?" || arg == "-h" || arg == "--help" {
      return errors.New(p.Info())
    } else if strings.HasPrefix(arg, LONG_OPT_PREFIX) {
      key := arg[2:]

      if strings.Contains(key, "=") {
        parts := strings.Split(key, "=")
        key = parts[0]
        if len(key) < 2 {
          return p.Error("expected at least two letters for long option")
        }

        if len(parts) > 2 {
          return p.Error("too many = in option " + arg)
        } else if len(parts) < 2 {
          panic("unexpected")
        }

        opt := p.findLongOption(key)
        if opt == nil {
          return p.Error("" + arg + " option not a valid option")
        }

        var optionArgs []string
        var err error
        if len(parts[1]) > 0 {
          optionArgs = strings.Split(parts[1], ",")
          if len(optionArgs) != opt.NArgs() {
            return p.Error("option --" + opt.Long() + " expects " + strconv.Itoa(opt.NArgs()) + ", got " + strconv.Itoa(len(optionArgs)))
          }
        } else if optionArgs, err = p.Eat(opt.NArgs()); err != nil {
          return err
        }

        if err := opt.Handle(optionArgs); err != nil {
          return p.Error(err.Error())
        }
      } else {
        if len(key) < 2 {
          return p.Error("expected at least two letters for long option")
        }

        opt := p.findLongOption(key)
        if opt == nil {
          return p.Error(arg + " option not valid")
        }

        optionArgs, err := p.Eat(opt.NArgs())
        if err != nil {
          return err
        }

        if err := opt.Handle(optionArgs); err != nil {
          return p.Error(err.Error())
        }
      }
    } else if strings.HasPrefix(arg, SHORT_OPT_PREFIX) {
      if len(arg) < 2 {
        return p.Error("plain dash not handled")
      }

      key := arg[1:2]

      opt := p.findShortOption(key)
      if opt == nil {
        return p.Error(arg + " option not valid")
      }

      if opt.NArgs() == 0 { // flag
        if err := opt.Handle([]string{}); err != nil {
          return p.Error(err.Error())
        }

        // remainer of key can be other flag
        for iChar, _ := range arg[2:] {
          addKey := arg[iChar:iChar+1]
          addOpt := p.findShortOption(addKey)
          if addOpt == nil {
            return p.Error(arg + " option not valid")
          }

          if addOpt.NArgs() != 0 {
            return p.Error(addKey + " is not a flag (part of " + arg + ")")
          }

          if err := addOpt.Handle([]string{}); err != nil {
            return p.Error(err.Error())
          }
        }
      } else {
        var optionArgs []string
        var err error
        if len(arg) > 2 {
          if arg[2] == '=' {
            rhs := arg[3:]
            if len(rhs) == 0 {
              optionArgs, err = p.Eat(opt.NArgs())
              if err != nil {
                return err
              }
            } else {
              optionArgs = strings.Split(rhs, ",")
            }
          } else {
            optionArgs = []string{arg[2:]}

            if opt.NArgs() > 1 {
              moreArgs, err := p.Eat(opt.NArgs()-1)
              if err != nil {
                return err
              }

              optionArgs = append(optionArgs, moreArgs...)
            }
          }
        } else if optionArgs, err = p.Eat(opt.NArgs()); err != nil {
          return err
        }

        if len(optionArgs) != opt.NArgs() {
          return p.Error("expected " + strconv.Itoa(opt.NArgs()) + " comma separated arguments after " + key + "=")
        }

        if err := opt.Handle(optionArgs); err != nil {
          return p.Error(err.Error())
        }
      }
    } else if runtime.GOOS == "windows" && strings.HasPrefix(arg, WIN_OPT_PREFIX) {
      key := arg[1:]

      if len(key) == 1 {
        opt := p.findShortOption(key)
        if opt == nil {
          return p.Error("invalid option " + arg)
        }

        optionArgs, err := p.Eat(opt.NArgs())
        if err != nil {
          return err
        }

        if err := opt.Handle(optionArgs); err != nil {
          return p.Error(err.Error())
        }
      } else if len(key) > 1 {
        opt := p.findLongOption(key)
        if opt == nil {
          return p.Error("invalid option " + arg)
        }

        optionArgs, err := p.Eat(opt.NArgs())
        if err != nil {
          return err
        }

        if err := opt.Handle(optionArgs); err != nil {
          return p.Error(err.Error())
        }
      } else {
        return p.Error("invalid option " + arg)
      }
    } else {
      positional = append(positional, arg)
    }
  }

  if p.last == nil {
    if len(positional) != 0 {
      return p.Error("expected 0 positional args")
    } else {
      return nil
    }
  } 

  if len(positional) != p.last.NArgs() {
    return p.Error("expected " + strconv.Itoa(p.last.NArgs()) + " positional args, got " + strconv.Itoa(len(positional)))
  }

  return p.last.Handle(positional)
}
