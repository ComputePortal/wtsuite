package scripts

import (
	"errors"
	"fmt"
  "sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/macros"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
)

type FileBundle struct {
	cmdDefines map[string]string
	scripts    []FileScript
}

func NewFileBundle(cmdDefines map[string]string) *FileBundle {
	return &FileBundle{cmdDefines, make([]FileScript, 0)}
}

func (b *FileBundle) newScope() js.GlobalScope {
	return &FileBundleScope{js.NewFilledGlobalScope(), b}
}

func (b *FileBundle) Append(s FileScript) {
	b.scripts = append(b.scripts, s)
}

func (b *FileBundle) IsEmpty() bool {
	return len(b.scripts) == 0
}

func (b *FileBundle) Write() (string, error) {
	var sb strings.Builder

	sb.WriteString(js.WriteGlobalHeaders())
	sb.WriteString(macros.WriteHeaders())

	for k, defineVal := range b.cmdDefines {
		sb.WriteString("const ")
		sb.WriteString(k)
		sb.WriteString("=\"")
		sb.WriteString(defineVal)
		sb.WriteString("\";")
	}

	for _, s := range b.scripts {
		str, err := s.Write()
		if err != nil {
			return sb.String(), err
		}

		if VERBOSITY >= 2 {
			fmt.Printf("%s\n", files.Abbreviate(s.Path()))
		}

		sb.WriteString(str)
	}

	return sb.String(), nil
}

// for jspp library functionality
func (b *FileBundle) resolveDependencies(s FileScript, deps *map[string]FileScript) error {
	callerCtx := s.Module().Context()
	callerPath := callerCtx.Path()

	for _, d := range s.Dependencies() {
		files.AddCacheDependency(callerPath, d)

		if _, ok := (*deps)[d]; !ok {
			new, err := NewFileScript(d, callerPath)
			if err != nil {
				if err.Error() == "not found" {
					errCtx := s.Module().Context()
					return errCtx.NewError("Error: '" + d + "' not found (from '" + callerPath + "')")
				} else {
					return err
				}
			}
			(*deps)[d] = new
			if err := b.resolveDependencies(new, deps); err != nil {
				return err
			}
		}

	}

	return nil
}

func (b *FileBundle) reportCircularDependencyRecursive(downstream []FileScript, fs FileScript, deps map[string]FileScript) error {
	for _, ds := range downstream {
		if ds.Path() == fs.Path() {
			return errors.New("Circular dependency found:\n")
		}
	}

	for _, d := range fs.Dependencies() {
		if err := b.reportCircularDependencyRecursive(append(downstream, fs), deps[d], deps); err != nil {
			return errors.New(err.Error() + " -> " + files.Abbreviate(deps[d].Path()) + "\n")
		}
	}

	return nil
}

func (b *FileBundle) reportCircularDependency(start FileScript, deps map[string]FileScript) error {
	for _, d := range start.Dependencies() {
		if err := b.reportCircularDependencyRecursive([]FileScript{start}, deps[d], deps); err != nil {
			return errors.New(err.Error() + " -> " + files.Abbreviate(deps[d].Path()) + "\n")
		}
	}

	return nil
}

// block recursion
func (b *FileBundle) ResolveDependencies() error {
  // first sort the already collected scripts alphabetically by path
  fss := NewFileScriptSorter(b.scripts)
  sort.Stable(fss)
  scripts := fss.Result()

	deps := make(map[string]FileScript)

	sortedScripts := make([]FileScript, 0)
	doneScripts := make(map[string]FileScript)
	unsortedScripts := make([]FileScript, 0)

	allDone := func(fs FileScript) bool {
		ok := true
		for _, d := range fs.Dependencies() {
			if _, ok_ := doneScripts[d]; !ok_ {
				ok = false
				break
			}
		}

		return ok
	}

	addToDone := func(fs FileScript) {
		if _, ok := doneScripts[fs.Path()]; !ok {
			sortedScripts = append(sortedScripts, fs)
			doneScripts[fs.Path()] = fs
		}
	}


	for _, s := range scripts {
		if err := b.resolveDependencies(s, &deps); err != nil {
			return err
		}

		if allDone(s) {
			addToDone(s)
		} else {
			unsortedScripts = append(unsortedScripts, s)
		}
	}

  depsKeys := make([]string, 0)
  for k, _ := range deps {
    depsKeys = append(depsKeys, k)
  }
  sort.Strings(depsKeys)

  for _, k := range depsKeys {
    fs := deps[k]
		if allDone(fs) {
			addToDone(fs)
		} else {
			unsortedScripts = append(unsortedScripts, fs)
		}
  }

	for len(unsortedScripts) > 0 {
		prevUnsortedScripts := unsortedScripts
		unsortedScripts = make([]FileScript, 0)

		for _, fs := range prevUnsortedScripts {
			if allDone(fs) {
				addToDone(fs)
			} else {
				unsortedScripts = append(unsortedScripts, fs)
			}
		}

		if len(unsortedScripts) > 0 && len(unsortedScripts) == len(prevUnsortedScripts) {
			// report circular dependency, which can start from any of the scripts
			err := b.reportCircularDependency(unsortedScripts[0], deps)
			if err == nil {
				panic("unable to find circular dep, but it must be there")
			}

			return err
		}
	}

	b.scripts = make([]FileScript, 0)
	for _, s := range sortedScripts {
		b.scripts = append(b.scripts, s)
	}

  if (VERBOSITY >= 2) {
    for _, s := range b.scripts {
      fmt.Printf("dep: %s\n", files.Abbreviate(s.Path()))
    }
  }

	return nil
}

func (b *FileBundle) ResolveNames() error {
	bs := b.newScope()

	for k, str := range b.cmdDefines {
		if bs.HasVariable(k) {
			return errors.New("Error: cmd define " + k + " already defined elsewhere")
		}

    ctx := context.NewDummyContext()
		variable := js.NewVariable(k, true, ctx)
    variable.SetValue(prototypes.NewLiteralString(str, ctx))

		// dont bother renaming, so we dont need to keep the newly created variable
		if err := bs.SetVariable(k, variable); err != nil {
			panic(err)
		}
	}

	for _, s := range b.scripts {
		if err := s.ResolveNames(bs); err != nil {
			return err
		}
	}

	return nil
}

func (b *FileBundle) EvalTypes() error {
	for _, s := range b.scripts {
		if err := s.EvalTypes(); err != nil {
			return err
		}
	}

	return nil
}

func (b *FileBundle) ResolveActivity() error {
	usage := js.NewUsage()

	// reverse stack order!
	for i := len(b.scripts) - 1; i >= 0; i-- {
		s := b.scripts[i]
		if err := s.ResolveActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (b *FileBundle) UniqueNames() error {
	ns := js.NewNamespace(nil, false)

	for _, s := range b.scripts {
		if err := s.UniqueEntryPointNames(ns); err != nil {
			return err
		}
	}

	for _, s := range b.scripts {
		if err := s.UniversalNames(ns); err != nil {
			return err
		}
	}

	for _, s := range b.scripts {
		if err := s.UniqueNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (b *FileBundle) Walk(fn func(scriptPath string, obj interface{}) error) error {
  for _, s := range b.scripts {
    fmt.Println(s.Path())
    if err := s.Walk(fn); err != nil {
      return err
    }
  }

  return nil
}

func (b *FileBundle) Finalize() error {
	if err := b.ResolveDependencies(); err != nil {
		return err
	}

	if err := b.ResolveNames(); err != nil {
		return err
	}

	if err := b.EvalTypes(); err != nil {
		return err
	}

	if err := b.ResolveActivity(); err != nil {
		return err
	}

	if err := b.UniqueNames(); err != nil {
		return err
	}

	return nil
}
