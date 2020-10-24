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

	"../../pkg/cache"
	"../../pkg/directives"
	"../../pkg/files"
	"../../pkg/parsers"
	"../../pkg/tokens/context"
	"../../pkg/tokens/js"
	"../../pkg/tokens/js/prototypes"
	"../../pkg/tokens/js/values"
	"../../pkg/tokens/html"
	"../../pkg/tree/scripts"
)

type CmdArgs struct {
  includeDirs []string
  operation string // eg. rename-package 

  // TODO: require file in case of ambiguity
  dryRun  bool
  verbosity int
  args []string // remaining positional args
}

func printUsageAndExit() {
  fmt.Fprintf(os.Stderr, "Usage: %s [options] --operation <operation> <args>\n", os.Args[0])

	fmt.Fprintf(os.Stderr, `
Options:
  -?, -h, --help         Show this message, other options are ignored
  -I, --include <dir>    Append a search directory to HTMLPPPATH
  -n                     Dry run
  -v[v[v[v...]]]         Verbosity

Operations:
  rename-package <old-name> <new-name>    Change package name, move its directory
  rename-class <old-name> <new-name>      Change class name, move its file
`)

  os.Exit(1)
}

func printMessageAndExit(msg string, printUsage bool) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")

	if printUsage {
		printUsageAndExit()
	} else {
		os.Exit(1)
	}
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
		includeDirs:     make([]string, 0),
    operation: "",
    dryRun: false,
		verbosity:     0,
	}

	positional := make([]string, 0)

	i := 1
	n := len(os.Args)

	for i < n {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-v") {
			re := regexp.MustCompile(`^[\-][v]+$`)
			if !re.MatchString(arg) {
				printMessageAndExit("Error: bad verbosity option "+arg, true)
			} else {
				cmdArgs.verbosity += len(arg) - 1
			}
		} else if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-?", "-h", "--help":
				printUsageAndExit()
			case "-I", "--include":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					cmdArgs.includeDirs = append(cmdArgs.includeDirs, os.Args[i+1])
					i++
				}
      case "--operation":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.operation != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.operation = os.Args[i+1]
					i++
				}
      case "-n":
        cmdArgs.dryRun = true
			}
		} else {
			positional = append(positional, arg)
		}

		i++
	}


	for _, includeDir := range cmdArgs.includeDirs {
		if err := files.AssertDir(includeDir); err != nil {
			printMessageAndExit("Error: include dir '"+includeDir+"' "+err.Error(), false)
		}
	}

  switch cmdArgs.operation {
  case "":
    printMessageAndExit("Error: --operation must be specified", true)
  case "rename-package":
    if len(positional) != 2 {
      printMessageAndExit("Error: both --old-name and --new-name must be set for rename-package operation", true)
    }
  case "rename-class":
    if len(positional) != 2 {
      printMessageAndExit("Error: both --old-name and --new-name must be set for rename-class operation", true)
    }
  default:
    printMessageAndExit("Error: unrecognized --operation " + cmdArgs.operation, true)
  }

  cmdArgs.args = positional

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
	files.JS_MODE = true

	js.TARGET = "all"
	directives.ForceNewViewFileScriptRegistration()
  directives.IGNORE_UNSET_URLS = true
  js.ALLOW_DUMMY_VIEW_INTERFACE = true

  html.PX_PER_REM = 16
	cache.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	prototypes.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

	files.AppendIncludeDirs(cmdArgs.includeDirs)
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
  case "rename-package":
    return renamePackage(bundle, cmdArgs.dryRun, cmdArgs.args[0], cmdArgs.args[1])
  case "rename-class":
    return renameClass(bundle, cmdArgs.dryRun, cmdArgs.args[0], cmdArgs.args[1])
  default:
    panic("not yet implemented")
  }
}

func renamePackage(bundle *scripts.FileBundle, dryRun bool, oldName string, newName string) error {
  // TODO: handle ambiguous packages
  // first we must find the package
  pwd, err := os.Getwd()
  if err != nil {
    return err
  }

  pkgPath, isPackage, err := files.SearchPackage(pwd, oldName, files.JSPACKAGE_SUFFIX)
  if err != nil {
    return err
  } else if !isPackage {
    return errors.New("Error: " + oldName + " is not a package")
  }

  pkgPath = strings.TrimSpace(pkgPath)
  pkgPathDir := filepath.Dir(pkgPath)
  newPkgPathDir := filepath.Join(pkgPathDir, filepath.Join("..", newName))

  fmt.Fprintf(os.Stdout, "Found package %s\n", pkgPath)

  contexts := make([]context.Context, 0)
  // we can now walk the syntax tree to detect an VarExpression that refers to this package
  if err := bundle.Walk(func(scriptPath string, obj_ interface{}) error {
    switch obj := obj_.(type) {
    case *js.VarExpression:
      // detect if part of package
      if obj.PackagePath() == pkgPath {
        ctx := obj.PackageContext()
        // this context should be pretty well refined
        contexts = append(contexts, ctx)
      }
    case *js.Member:
      if obj.PackagePath() == pkgPath {
        ctx := obj.ObjectContext()
        // this context should be pretty well refined
        contexts = append(contexts, ctx)
      }
    case *js.ImportedVariable:
      objImportAbsPath := obj.AbsPath()
      if objImportAbsPath == pkgPath {
        ctx := obj.PathContext()
        
        origPath := ctx.Content()

        origPathParts := strings.Split(filepath.ToSlash(origPath), "/")

        lastNonMatch := 0
        for i, part := range origPathParts {
          if part == oldName && i > 0 { // XXX: is this test good enough, or should we check full path?
            lastNonMatch += len(origPathParts[i-1]) + 1
          }
        }

        if lastNonMatch > 0 {
          ctx = ctx.NewContext(lastNonMatch, len(origPath))
        }

        contexts = append(contexts, ctx)
      } else if strings.HasPrefix(objImportAbsPath, pkgPathDir) {
        // the end that is not part cannot be replaced
        ctx := obj.PathContext()

        origPath := ctx.Content()

        // cut this from end of the complete path
        fullPathRoot := filepath.Dir(ctx.Path())

        origPathParts := strings.Split(filepath.ToSlash(origPath), "/")

        found := false
        start := 0
        end := len(origPath)

        tmp := fullPathRoot
        for _, part := range origPathParts {
          tmp = filepath.Join(tmp, part)
          if filepath.Clean(tmp) == filepath.Clean(pkgPathDir) {
            found = true
            end = start + len(part)
            break
          } 

          start += len(part) + 1
        }

        if found {
          // only create the newContext if exactly the package name is found in the origPath, and it refers the correct directory
          ctx = ctx.NewContext(start, end)

          if ctx.Content() != "." {
            contexts = append(contexts, ctx)
          }
        } 
      }
    }

    return nil
  }); err != nil {
    return err
  }

  moveMap := make(map[string]string)
  moveMap[pkgPathDir] = newPkgPathDir

  if err := renameSymbolsAndMoveFiles(dryRun, contexts, oldName, newName, moveMap); err != nil {
    return err
  }

  return nil
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

  setUpEnv(cmdArgs)

  // setup the cache, even though it isn't needed (to even nil pointer derefence errors in some places
	cache.LoadJSCache("", true)

  if err := applyOperation(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }
}
