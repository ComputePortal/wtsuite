// TODO: 
// * rename-instance
// * rename-function

package main

import (
  "errors"
  "fmt"
  "os"
  "path/filepath"
  "regexp"
  "strings"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

const DEFAULT_OPERATION = "rename-class"

var cmdParser *parsers.CLIParser = nil

type CmdArgs struct {
  operation string // eg. rename-package 

  // TODO: require file in case of ambiguity
  dryRun  bool
  verbosity int
  args []string // remaining positional args
}

func printMessageAndExit(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
  os.Exit(1)
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
    operation: DEFAULT_OPERATION,
    dryRun: false,
		verbosity:     0,
	}

	var positional []string = nil

  cmdParser = parsers.NewCLIParser(fmt.Sprintf("Usage: %s [options] [--op-type <operation>] <args>\n", os.Args[0]), 
  `Operations:
  rename-class <old-name> <new-name>      Change class name, move its file`,
    []parsers.CLIOption{
      parsers.NewCLIUniqueFlag("n", ""       , "-n              Dry run", &(cmdArgs.dryRun)),
      parsers.NewCLIUniqueEnum("t", "type"   , "-t, --op-type   <operation-type>    Defaults to \"" + DEFAULT_OPERATION + "\", see below for other possibilities", []string{"rename-class"}, &(cmdArgs.operation)),
      parsers.NewCLIUniqueFlag("l", "latest" , "-l, --latest    Ignore max semver, use latest tagged versions of dependencies", &(files.LATEST)),
      parsers.NewCLICountFlag("v", ""        , "-v[v[v..]]      Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIRemaining(&positional),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

  switch cmdArgs.operation {
  case "rename-class":
    if len(positional) != 2 {
      printMessageAndExit("Error: both --old-name and --new-name must be set for rename-class operation")
    }
  default:
    panic("unhandled") // CLIUniqueEnum should've been able to catch this
  }

  cmdArgs.args = positional

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) error {
	files.JS_MODE = true

	js.TARGET = "all"
	directives.ForceNewViewFileScriptRegistration(directives.NewFileCache())
  directives.IGNORE_UNSET_URLS = true

  html.PX_PER_REM = 16
	cache.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

  pwd, err := os.Getwd()
  if err != nil {
    return err
  }

  return files.ResolvePackages(filepath.Join(pwd, files.PACKAGE_JSON))
}

func applyOperation(cmdArgs CmdArgs) error {
  // only add files once (abs path -> scripts.FileScript)
  bundle := scripts.NewFileBundle(map[string]string{})

  pwd, err := os.Getwd()
  if err != nil {
    return err
  }

  // TODO: cmdArg so we can walk different directory
  if err := files.WalkFiles(pwd, files.JSFILE_EXT, func(path string) error {
    // caller be left empty because path is absolute
    if !filepath.IsAbs(path) {
      panic(path + " should be absolute")
    }
    fs, err := scripts.NewFileScript(path, "")
    if err != nil {
      return err
    }

    bundle.Append(fs)

    return nil
  }); err != nil {
    return err
  }

  // all scripts should be included, but they need to be sorted
  if err := bundle.ResolveDependencies(); err != nil {
    return err
  }

  if err := bundle.ResolveNames(); err != nil {
    return err
  }

  switch cmdArgs.operation {
  case "rename-class":
    return renameClass(bundle, cmdArgs.dryRun, cmdArgs.args[0], cmdArgs.args[1])
  default:
    panic("not yet implemented")
  }
}

func renameClass(bundle *scripts.FileBundle, dryRun bool, oldName string, newName string) error {
  // walk a first time to find the class
  var class *js.Class = nil

  if err := bundle.Walk(func(_ string, obj_ interface{}) error {
    if obj, ok := (obj_).(*js.Class); ok {
      if obj.Name() == oldName {
        if (class == nil) {
          class = obj
        } else if (class != obj) {
          return errors.New("Error: class " + oldName + " is ambiguous")
        }
      }
    }
    return nil
  }); err != nil {
    return err
  }

  if class == nil {
    return errors.New("Error: class " + oldName + " not found")
  }

  // now find out if we must rename file containing the class
  classCtx := class.Context()
  filePath := classCtx.Path()
  ext := filepath.Ext(filePath)
  fileBaseName := strings.TrimRight(filepath.Base(filePath), ext)

  moveFileToo := fileBaseName == oldName

  fmt.Fprintf(os.Stdout, "Found class %s in %s\n", oldName, filePath)

  // now collect all the contexts
  contexts := make([]context.Context, 0)

  // only do import paths once, even though they might be used for several symbols
  donePathLiterals := make(map[interface{}]bool)

  if err := bundle.Walk(func(scriptPath string, obj_ interface{}) error {
    switch obj := obj_.(type) {
    case *js.VarExpression:
      refObj_ := obj.GetVariable().GetObject()
      if refObj_ != nil {
        if refObj, ok := refObj_.(*js.Class); ok {
          if refObj == class {
            ctx := obj.NonPackageContext()
            contexts = append(contexts, ctx)
          }
        }
      }
    case *js.Member:
      // the first condition is that the member must evaluate to the class
      _, keyValue := obj.ObjectNameAndKey()
      if keyValue == oldName {
        pkgMember, err := obj.GetPackageMember() 
        if err != nil {
          return err
        }

        refObj_ := pkgMember.GetObject()
        if refObj_ != nil {
          if refObj, ok := refObj_.(*js.Class); ok {
            if refObj == class {
              ctx := obj.KeyContext()
              contexts = append(contexts, ctx)
            }
          }
        }
      }
    case *js.ImportedVariable:
      // only exact match is possible because these cannot be directories
      if moveFileToo && obj.AbsPath() == filePath {
        if _, ok := donePathLiterals[obj.PathLiteral()]; !ok {
          ctx := obj.PathContext()

          origPath := ctx.Content()
          if strings.HasPrefix(origPath, "\"") {
            panic("can't start with quotes")
          }

          origDir := strings.TrimRight(origPath, filepath.Base(filePath))

          ctx = ctx.NewContext(len(origDir), len(origPath) - len(ext))
          contexts = append(contexts, ctx)

          donePathLiterals[obj.PathLiteral()] = true
        }
      }
      
      v := obj.GetVariable()
      if v != nil {
        refObj_ := v.GetObject()
        if refObj, ok := refObj_.(*js.Class); ok && refObj == class {
          ctx := obj.PathContext()
          origPath := ctx.Content()

          completeCtx := obj.Context()

          completeContent := completeCtx.Content()

          // cut off the path part, which is always last for static imports (dynamic imports dont have any reference to the class anyway
          completeContent = strings.TrimRight(completeContent, origPath)

          // if newName==oldName and it happens to be a part of this content, then it is also replaced
          if strings.Contains(completeContent, oldName) {
            // special regexp should be used to only replace up till word boundaries
            re := regexp.MustCompile(`\b` + oldName + `\b`)

            indices := re.FindAllStringIndex(completeContent, -1)

            if indices == nil {
              panic("unexpected due to contains check")
            }

            for _, idx := range indices {
              extraCtx := completeCtx.NewContext(idx[0], idx[1])
              contexts = append(contexts, extraCtx)
            }
          }
        }
      }
    }

    return nil
  }); err != nil {
    return err
  }

  moveMap := make(map[string]string)
  if moveFileToo {
    moveMap[filePath] = filepath.Join(filepath.Dir(filePath), newName + ext)
  }

  if err := renameSymbolsAndMoveFiles(dryRun, contexts, 
    oldName, newName, moveMap); err != nil {
    return err
  }

  return nil
}

// rename contexts are merged
// move must come after symbol renaming
// files to be moved can also be directories
// XXX: hopefully no errors occur here, because then the files will be mangled
func renameSymbolsAndMoveFiles(dryRun bool, contexts []context.Context, 
  oldName, newName string, moveMap map[string]string) error {
  
  // check that the move is possible (newFNames cant exists)
  for oldFName, newFName := range moveMap {
    if _, err := os.Stat(newFName); !os.IsNotExist(err) {
      return errors.New("Error: can't move " + oldFName + ", " + newFName + " already exists")
    }
  }

  if dryRun {
    fmt.Fprintf(os.Stdout, "#Found %d symbols, and %d files to rename\n", len(contexts), len(moveMap))
    // print the contexts nicely
    for _, ctx := range contexts {
      fmt.Fprintf(os.Stdout, ctx.WritePrettyOneLiner())
    }

    for oldFile, newFile := range moveMap {
      fmt.Fprintf(os.Stdout, "\u001b[35m%s\u001b[0m -> \u001b[35m%s\u001b[0m\n", oldFile, newFile)
    }
  } else {
    // contexts on the same file must be merged
    fileContexts := make(map[string]context.Context)

    for _, ctx := range contexts {
      p := ctx.Path()

      if prevCtx, ok := fileContexts[p]; ok {
        fileContexts[p] = prevCtx.Merge(ctx)
      } else {
        fileContexts[p] = ctx
      }
    }

    for _, ctx := range fileContexts {
      if err := ctx.SearchReplaceOrig(oldName, newName); err != nil {
        return err
      }
    }

    // must come after SearchReplaceOrig because contexts use original filenames to write new symbol names
    for oldFile, newFile := range moveMap {
      if err := os.Rename(oldFile, newFile); err != nil {
        return err
      }
    }

    fmt.Fprintf(os.Stdout, "#Renamed %d locations to rename and moved %d file\n", 
      len(contexts), len(moveMap))
  }

  return nil
}

func main() {
  cmdArgs := parseArgs()

  if err := setUpEnv(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }

  // setup the cache, even though it isn't needed (to even nil pointer derefence errors in some places
	cache.LoadJSCache("", true)

  if err := applyOperation(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }
}
