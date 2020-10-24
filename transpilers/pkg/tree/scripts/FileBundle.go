package scripts

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"../../files"
	"../../tokens/context"
	"../../tokens/js"
	"../../tokens/js/macros"
	"../../tokens/js/prototypes"
)

const NO_EVAL = false // true for debuging

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

	for _, s := range b.scripts {
		if err := b.resolveDependencies(s, &deps); err != nil {
			return err
		}

		if allDone(s) {
			addToDone(s)
		} else {
			unsortedScripts = append(unsortedScripts, s)
		}
	}

	for _, fs := range deps {
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

	return nil
}

func (b *FileBundle) ResolveNames() error {
	bs := b.newScope()

	for _, s := range b.scripts {
		if err := s.ResolveNames(bs); err != nil {
			return err
		}
	}

	return nil
}

func (b *FileBundle) ResolveControlNames(viewInterfaces map[string]*js.ViewInterface) (map[string]js.Variable, error) {
	bs := b.newScope()

	for _, vif := range viewInterfaces {
		if err := vif.ResolveNames(bs); err != nil {
			panic(err)
		}
	}

	cmdDefineVars := make(map[string]js.Variable)
	for k, _ := range b.cmdDefines {
		if bs.HasVariable(k) {
			return nil, errors.New("Error: cmd define " + k + " already defined elsewhere")
		}

		v := js.NewVariable(k, true, context.NewDummyContext())
		cmdDefineVars[k] = v

		// dont bother renaming, so we dont need to keep the newly created variable
		if err := bs.SetVariable(k, v); err != nil {
			panic(err)
		}
	}

	for _, s := range b.scripts {
		if err := s.ResolveNames(bs); err != nil {
			return nil, err
		}
	}

	return cmdDefineVars, nil
}

func (b *FileBundle) EvalTypes(stack *js.GlobalStack) error {
	for _, s := range b.scripts {
		if err := s.EvalTypes(stack); err != nil {
			return err
		}
	}

	return nil
}

//func (b *FileBundle) evalControlTypes(globalStack *js.GlobalStack,
//viewInterfaces map[string]*js.ViewInterface) error {
//}

func (b *FileBundle) CreateNewGlobalStack(cacheStack *js.CacheStack, cmdDefineVars map[string]js.Variable) *js.GlobalStack {
	globalStack := js.NewFilledGlobalStack(cacheStack)

	for k, defineVal := range b.cmdDefines {
		defineVar := cmdDefineVars[k]

		ctx := context.NewDummyContext()
		if err := globalStack.SetValue(defineVar, prototypes.NewLiteralString(defineVal, ctx), false, ctx); err != nil {
			panic(err)
		}
	}

	return globalStack
}

type empty struct{}

// use the list of updatedViews and updatedControls to determine which controls actually need to be evaluated
func (b *FileBundle) EvalControlTypes(viewInterfaces map[string]*js.ViewInterface, cmdDefineVars map[string]js.Variable, updatedViews []string, updatedControls []string) error {
	// for each control script, for each related (unique) viewInterface: setup a ViewInterfaceStack, and eval all the types
	fns := make([]func(string) error, 0)
	fnViews := make([]string, 0)

	cacheStack := js.NewCacheStack()
	globalStack_ := b.CreateNewGlobalStack(cacheStack, cmdDefineVars)
	for _, s_ := range b.scripts {
		if s, ok := s_.(*ControlFileScript); ok {
			isUpdatedControl := false
			for _, uc := range updatedControls {
				if uc == s.Path() {
					isUpdatedControl = true
					break
				}
			}

			prevViewInterfaces := make([]*js.ViewInterface, 0)

			for _, view := range s.views {
				isUpdatedView := false
				for _, uv := range updatedViews {
					if uv == view {
						isUpdatedView = true
						break
					}
				}

				if !(isUpdatedView || isUpdatedControl) {
					continue
				}

				if VERBOSITY >= 3 {
					fmt.Fprintf(os.Stderr, "Evaluating for view interface %s\n", view)
				}

				viewInterface := viewInterfaces[view]

				unique := true
				for _, prev := range prevViewInterfaces {
					if prev.IsSame(viewInterface) {
						unique = false
					}
				}

				if unique {
					fn := func(view string) error {
						var globalStack *js.GlobalStack
						if js.SERIAL {
							globalStack = globalStack_
						} else {
							globalStack = b.CreateNewGlobalStack(cacheStack, cmdDefineVars)
						}
						// can we parallelize this?
						stack := js.NewViewInterfaceStack(viewInterfaces[view], viewInterfaces, globalStack)
						for _, sInner_ := range b.scripts {
							if sInner, ok := sInner_.(*ControlFileScript); ok && !sInner.hasView(view) {
								continue
							}

							if err := sInner_.EvalTypes(stack); err != nil {
								return err
							}
						}

						return nil
					}

					fnViews = append(fnViews, view)
					fns = append(fns, fn)
					//if err := fn(); err != nil {
					//return err
					//}

					prevViewInterfaces = append(prevViewInterfaces, viewInterface)
				}
			}
		}
	}

	// actually evaluate the pending functions
	if js.SERIAL {
		for i, fn := range fns {
			if err := fn(fnViews[i]); err != nil {
				return err
			}
		}
	} else {

		sem := make(chan empty, len(fns))
		errs := make([]error, len(fns))

		for i, fn := range fns {
			go func(i int, view string) {
				errs[i] = fn(view)
				sem <- empty{}
			}(i, fnViews[i])
		}

		for i := 0; i < len(fns); i++ {
			<-sem
		}

		for _, err := range errs {
			if err != nil {
				return err
			}
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

	stack := js.NewFilledGlobalStack(js.NewCacheStack())

	if err := b.EvalTypes(stack); err != nil {
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

func (b *FileBundle) FinalizeControls(viewInterfaces map[string]*js.ViewInterface, updatedViews []string, updatedControls []string) error {
	if err := b.ResolveDependencies(); err != nil {
		return err
	}

	if VERBOSITY >= 2 {
		fmt.Println("Done resolving dependencies")
	}

	cmdDefineVars, err := b.ResolveControlNames(viewInterfaces)
	if err != nil {
		return err
	}

	if VERBOSITY >= 2 {
		fmt.Println("Done control names")
	}

	if !NO_EVAL {
		if err := b.EvalControlTypes(viewInterfaces, cmdDefineVars, updatedViews, updatedControls); err != nil {
			return err
		}
	}

	if VERBOSITY >= 2 {
		fmt.Println("Done evaluating type")
	}

	// viewInterfaces use builtin types so dont require 'ResolveActivity()' or 'UniqueNames()' step

	if err := b.ResolveActivity(); err != nil {
		return err
	}

	if err := b.UniqueNames(); err != nil {
		return err
	}

	return nil
}

