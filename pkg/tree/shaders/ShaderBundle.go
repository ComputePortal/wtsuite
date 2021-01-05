package shaders

import (
  "errors"
  "fmt"
  "sort"
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	jsv "github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type ShaderBundle struct {
  usage glsl.Usage
  version *glsl.Word // eg. "100 es"
  shaders []ShaderFile
}

func NewShaderBundle() *ShaderBundle {
  return &ShaderBundle{nil, nil, make([]ShaderFile, 0)}
}

func (b *ShaderBundle) Append(s ShaderFile) {
  b.shaders = append(b.shaders, s)
}

func (b *ShaderBundle) Write(nl string, tab string) (string, error) {
  var sb strings.Builder

  if b.version != nil {
    sb.WriteString("#version ")
    sb.WriteString(b.version.Value())
    sb.WriteString("\n")
  }

  for _, s := range b.shaders {
    str, err := s.Write(b.usage, nl, tab)
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

func (b *ShaderBundle) ResolveDependencies() error {
  // first sort the already collected shaders alphabetically by path
  fss := NewShaderFileSorter(b.shaders)
  sort.Stable(fss)
  shaders := fss.Result()

	deps := make(map[string]ShaderFile)

	sortedScripts := make([]ShaderFile, 0)
	doneScripts := make(map[string]ShaderFile)
	unsortedScripts := make([]ShaderFile, 0)

	allDone := func(fs ShaderFile) bool {
		ok := true
		for _, pl := range fs.Dependencies() {
      d := pl.Path
			if _, ok_ := doneScripts[d]; !ok_ {
				ok = false
				break
			}
		}

		return ok
	}

	addToDone := func(fs ShaderFile) {
		if _, ok := doneScripts[fs.Path()]; !ok {
			sortedScripts = append(sortedScripts, fs)
			doneScripts[fs.Path()] = fs
		}
	}


	for _, s := range shaders {
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
		unsortedScripts = make([]ShaderFile, 0)

		for _, fs := range prevUnsortedScripts {
			if allDone(fs) {
				addToDone(fs)
			} else {
				unsortedScripts = append(unsortedScripts, fs)
			}
		}

		if len(unsortedScripts) > 0 && len(unsortedScripts) == len(prevUnsortedScripts) {
			// report circular dependency, which can start from any of the shaders
			err := b.reportCircularDependency(unsortedScripts[0], deps)
			if err == nil {
				panic("unable to find circular dep, but it must be there")
			}

			return err
		}
	}

	b.shaders = make([]ShaderFile, 0)
	for _, s := range sortedScripts {
		b.shaders = append(b.shaders, s)
	}

  if (VERBOSITY >= 2) {
    for _, s := range b.shaders {
      fmt.Printf("dep: %s\n", files.Abbreviate(s.Path()))
    }
  }

	return nil
}

// TODO: choose between vertex or fragment target
func (b *ShaderBundle) newScope() glsl.GlobalScope {
	return &ShaderBundleScope{glsl.NewFilledGlobalScope(), b}
}

func (b *ShaderBundle) ResolveNames() error {
	bs := b.newScope()

	for _, s := range b.shaders {
		if err := s.ResolveNames(bs); err != nil {
			return err
		}
	}

	return nil
}

func (b *ShaderBundle) EvalTypes() error {
  for _, s := range b.shaders {
    if err := s.EvalTypes(); err != nil {
      return err
    }
  }

  return nil
}

func (b *ShaderBundle) resolveDependencies(s ShaderFile, deps *map[string]ShaderFile) error {
	callerCtx := s.Module().Context()
	callerPath := callerCtx.Path()

	for _, pl := range s.Dependencies() {
    d := pl.Path
    if files.AddCacheDependency != nil {
      files.AddCacheDependency(callerPath, d)
    }

		if _, ok := (*deps)[d]; !ok {
			new, err := NewShaderFile(d)
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

func (b *ShaderBundle) reportCircularDependencyRecursive(downstream []ShaderFile, fs ShaderFile, deps map[string]ShaderFile) error {
	for _, ds := range downstream {
		if ds.Path() == fs.Path() {
			return errors.New("Circular dependency found:\n")
		}
	}

	for _, pl := range fs.Dependencies() {
    d := pl.Path
		if err := b.reportCircularDependencyRecursive(append(downstream, fs), deps[d], deps); err != nil {
			return errors.New(err.Error() + " -> " + files.Abbreviate(deps[d].Path()) + "\n")
		}
	}

	return nil
}

func (b *ShaderBundle) reportCircularDependency(start ShaderFile, deps map[string]ShaderFile) error {
	for _, pl := range start.Dependencies() {
    d := pl.Path
		if err := b.reportCircularDependencyRecursive([]ShaderFile{start}, deps[d], deps); err != nil {
			return errors.New(err.Error() + " -> " + files.Abbreviate(deps[d].Path()) + "\n")
		}
	}

	return nil
}

func (b *ShaderBundle) ResolveActivity() error {
  b.usage = glsl.NewUsage()

  for i := len(b.shaders) - 1; i >= 0; i-- {
    s := b.shaders[i]
    if err := s.ResolveActivity(b.usage); err != nil {
      return err
    }
  }

  for _, s := range b.shaders {
    if err := s.FinalizeInjected(b.usage); err != nil {
      return err
    }
  }

  return b.usage.DetectUnused()
}

func (b *ShaderBundle) UniqueNames() error {
  ns := glsl.NewNamespace(nil, false)

  for _, s := range b.shaders {
    if err := s.UniqueEntryPointNames(ns); err != nil {
      return err
    }
  }

	for _, s := range b.shaders {
		if err := s.UniqueNames(ns); err != nil {
			return err
		}
	}

  return nil
}

func (b *ShaderBundle) CollectVersion() error {
  for _, s := range b.shaders {
    var err error
    if b.version, err = s.CollectVersion(b.version); err != nil {
      return err
    }
  }

  return nil
}

func (b *ShaderBundle) Finalize() error {
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

  if err := b.CollectVersion(); err != nil {
    return err
  }

  return nil
}

func (b *ShaderBundle) CollectVaryings(varyings map[string]string) error {
  for _, s := range b.shaders {
    if err := s.CollectVaryings(varyings); err != nil {
      return err
    }
  }

  return nil
}

func (b *ShaderBundle) FindExportedConst(name string) *glsl.Const {
  for _, s := range b.shaders {
    if cSt := s.FindExportedConst(name); cSt != nil {
      return cSt
    }
  }

  return nil
}

func (b *ShaderBundle) InjectConsts(rtName string, consts map[string]jsv.Value, ctx context.Context) error {
  for name, val := range consts {
    if jsv.IsAny(val) {
      errCtx := val.Context()
      return errCtx.NewError("Error: don't know how inject " + name + " of type " + val.TypeName() + " into a shader const")
    }

    cStatement := b.FindExportedConst(name)
    if cStatement == nil {
      return ctx.NewError("Error: exported const \"" + name + "\" not found in shader")
    }

    n := cStatement.Length()
    t := cStatement.TypeName()

    switch {
    case t == "int" && n < 2:
      if prototypes.IsInt(val) {
        cStatement.SetAltRHS("${" + rtName + "['" + name + "'].toString()}")
        continue
      }
    case t == "float" && n < 2:
      if prototypes.IsNumber(val) {
        cStatement.SetAltRHS("float(${" + rtName + "['" + name + "'].toString()})")
        continue
      }
    case t == "bool" && n < 2: 
      if prototypes.IsBoolean(val) {
        cStatement.SetAltRHS("${" + rtName + "['" + name + "'].toString())}")
        continue
      }
    }

    debugTypeName := t
    if n > 1 {
      debugTypeName += "[" + strconv.Itoa(n) + "]"
    }
    return ctx.NewError("Error: don't know how inject " + name + " of type " + val.TypeName() + " into a shader const of type " + debugTypeName)
  }

  return nil
}

// second return value is the map of all the varying types
func transpileWebGLShader(callerPath string, shaderPath_ *js.Word, rtName string, consts map[string]jsv.Value) (string, map[string]string, error) {
  errCtx := shaderPath_.Context()

  shaderPath, err := files.Search(callerPath, shaderPath_.Value())
  if err != nil {
    return "", nil, errCtx.NewError("Error: shader file \"" + shaderPath_.Value() + "\" not found")
  }

  bundle := NewShaderBundle()

  entryShader, err := NewInitShaderFile(shaderPath)
  if err != nil {
    return "", nil, errCtx.NewError("Error: problem reading shader file \"" + shaderPath_.Value() + "\" (" + err.Error() + ")")
  }

  bundle.Append(entryShader)

  if err := bundle.Finalize(); err != nil {
    return "", nil, err
  }

  if len(consts) > 0 {
    if err := bundle.InjectConsts(rtName, consts, errCtx); err != nil {
      return "", nil, err
    }
  }

  varyings := make(map[string]string)

  if err := bundle.CollectVaryings(varyings); err != nil {
    return "", nil, err
  }

  shaderSource, err := bundle.Write(patterns.NL, patterns.TAB)
  if err != nil {
    return "", nil, err
  }

  var b strings.Builder
  b.WriteString("`")
  b.WriteString(shaderSource)
  b.WriteString("`")

  return b.String(), varyings, nil
}

func TranspileWebGLShaders(callerPath string, vertexPath *js.Word, vertexConsts map[string]jsv.Value,
  fragmentPath *js.Word, fragmentConsts map[string]jsv.Value) (string, string, error) {

  vertexSource, vertexVaryings, err := transpileWebGLShader(callerPath, vertexPath, "v", vertexConsts)
  if err != nil {
    return "", "", err
  }

  fragmentSource, fragmentVaryings, err := transpileWebGLShader(callerPath, fragmentPath, "f", fragmentConsts)
  if err != nil {
    return "", "", err
  }

  errCtx := context.MergeContexts(vertexPath.Context(), fragmentPath.Context())

  for k, typeName := range vertexVaryings {
    fragTypeName, ok := fragmentVaryings[k] 
    if !ok {
      return "", "", errCtx.NewError("Error: varying " + k + " not found in fragment shader")
    }

    if fragTypeName != typeName {
      return "", "", errCtx.NewError("Error: varying " + k + " has different type in in fragment shader")
    }
  }

  for k, typeName := range fragmentVaryings {
    vertexTypeName, ok := vertexVaryings[k] 
    if !ok {
      return "", "", errCtx.NewError("Error: varying " + k + " not found in vertex shader")
    }

    if vertexTypeName != typeName {
      return "", "", errCtx.NewError("Error: varying " + k + " has different type in in vertex shader")
    }
  }

  if len(vertexVaryings) != len(fragmentVaryings) {
    return "", "", errCtx.NewError("Error: varyings differ")
  }

  return vertexSource, fragmentSource, nil
}
