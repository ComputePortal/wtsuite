package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "regexp"
  "strings"

	"../../pkg/cache"
	"../../pkg/directives"
	"../../pkg/files"
	"../../pkg/parsers"
	"../../pkg/tokens/js"
	"../../pkg/tokens/js/values"
	"../../pkg/tokens/html"
	"../../pkg/tree/scripts"
)

type CmdArgs struct {
  includeDirs []string // just for search
  graphType string // eg. class
  outputFile string // file needed by the graphviz utility 'dot' to create the visual, must be specified
  entryFile string // this is the entry point
  analyzedFiles map[string]string // (analyzing everything available in the include dirs would be too messy)

  verbosity int
}

func printUsageAndExit() {
  fmt.Fprintf(os.Stderr, "Usage: %s [options] --type <graph-type> --output <output-files> <input-files>\n", os.Args[0])

  fmt.Fprintf(os.Stderr, `
Options:
  -?, -h, --help             Show this message, other options are ignored
	-I, --include <dir>        Append a search directory to HTMLPPPATH
	-v[v[v[v...]]]             Verbosity
  --type <graph-type>        Graph type (see below)
  -o, --output <output-file> Output file location for graphviz dot

Graph types:
  class                  Class inheritance, explicit interface implementation
  instance               Instance properties
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
    graphType: "class",
    outputFile: "",
    entryFile: "",
    analyzedFiles: make(map[string]string),
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
      case "--type":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					cmdArgs.graphType = os.Args[i+1]
					i++
				}
      case "-o", "--output":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.outputFile != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.outputFile = os.Args[i+1]
					i++
				}
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

  if cmdArgs.outputFile == "" {
    printMessageAndExit("Error: --output not specified", true)
  }

  var err error
  cmdArgs.outputFile, err = filepath.Abs(cmdArgs.outputFile)
  if err != nil {
    printMessageAndExit("Error: " +err.Error(), false)
  }

  switch cmdArgs.graphType {
  case "class":
    if len(positional) == 0 {
      printMessageAndExit("Error: graph type class requires at least one input file", true)
    }
  case "instance":
    if len(positional) == 0 {
      printMessageAndExit("Error: graph type instance requires at least one input file", true)
    }
  default:
    printMessageAndExit("Error: unrecognized --type " + cmdArgs.graphType, true)
  }

  orderedInputFiles := make([]string, 0)

  for _, arg := range positional {
    info, err := os.Stat(arg)
    if os.IsNotExist(err) {
      printMessageAndExit("Error: \"" + arg + "\" not found", true)
    }

    if info.IsDir() {
      // walk to find the files
      if err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
        if filepath.Ext(path) == files.JSFILE_EXT {
          absPath, err := filepath.Abs(path)
          if err != nil {
            return err
          }

          orderedInputFiles = append(orderedInputFiles, absPath)
        }

        return nil
      }); err != nil {
        printMessageAndExit("Error: " + err.Error(), false)
      }
    } else {
      absPath, err := filepath.Abs(arg)
      if err != nil {
        printMessageAndExit("Error: " + err.Error(), false)
      }

      orderedInputFiles = append(orderedInputFiles, absPath)
    }
  }

  cmdArgs.entryFile = orderedInputFiles[0]

  // analyzedFiles also includes entryFile
  for _, path := range orderedInputFiles {
    cmdArgs.analyzedFiles[path] = path
  }

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
	files.JS_MODE = true

	js.TARGET = "all"
	directives.ForceNewViewFileScriptRegistration()
  directives.IGNORE_UNSET_URLS = true

  html.PX_PER_REM = 16
	cache.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

	files.AppendIncludeDirs(cmdArgs.includeDirs)
}

func createGraph(cmdArgs CmdArgs) error {
  // only add files once (abs path -> scripts.FileScript)
  bundle := scripts.NewFileBundle(map[string]string{})

  // also contains entryFile
  for _, scriptPath := range cmdArgs.analyzedFiles {
    fs, err := scripts.NewFileScript(scriptPath, "")
    if err != nil {
      return err
    }
    bundle.Append(fs)
  }

  if err := bundle.ResolveDependencies(); err != nil {
    return err
  }

  if err := bundle.ResolveNames(); err != nil {
    return err
  }

  // TODO: refactor graphing methods once we know how to tackle function dependencies
  switch cmdArgs.graphType {
  case "class":
    return createClassGraph(bundle, cmdArgs.entryFile, cmdArgs.analyzedFiles, cmdArgs.outputFile)
  case "instance":
    return createInstanceGraph(bundle, cmdArgs.entryFile, cmdArgs.analyzedFiles, cmdArgs.outputFile)
  default:
    panic("not yet implemented")
  }
}

func createClassGraph(bundle *scripts.FileBundle, entryFile string, 
  analyzedFiles map[string]string, outputFile string) error {
  var graph *Graph
  if len(analyzedFiles) == 1 {
    // only entry file
    graph = NewGraph(nil)
  } else {
    graph = NewGraph(analyzedFiles)
  }

  if err := bundle.Walk(func(scriptPath string, obj_ interface{}) error {
    if scriptPath != entryFile {
      // skip
      return nil
    }

    switch obj := obj_.(type) {
    case *js.Class:
      if err := graph.AddClass(obj); err != nil {
        return err
      }
    }

    return nil
  }); err != nil {
    return err
  }

  return writeGraph(graph, outputFile)
}

func createInstanceGraph(bundle *scripts.FileBundle, entryFile string, 
  analyzedFiles map[string]string, outputFile string) error {
  // needed so nodejs imports are set right
  if err := bundle.EvalTypes(); err != nil {
    return err
  }

  var graph *Graph
  if len(analyzedFiles) < 2 {
    graph = NewGraph(nil)
  } else {
    graph = NewGraph(analyzedFiles)
  }

  if err := bundle.Walk(func(scriptPath string, obj_ interface{}) error {
    if scriptPath != entryFile {
      // skip
      return nil
    }

    switch obj := obj_.(type) {
    case *js.Class:
      // try to instaniate the class (only instantiable classes are added)
      classVal, err := obj.GetClassValue()
      if err == nil {
        instance_, err := classVal.EvalConstructor(nil, obj.Context())
        if err == nil {
          if instance, ok := instance_.(*values.Instance); ok {
            if err := graph.AddInstance(instance); err != nil { // the used name will be 
              return err
            }
          }
        }
      }
    }

    return nil
  }); err != nil {
    return err
  }

  return writeGraph(graph, outputFile)
}
  
func writeGraph(graph *Graph, outputFile string) error {
  graph.Clean()

  result := graph.Write()

  if err := ioutil.WriteFile(outputFile, []byte(result), 0644); err != nil {
    return err
  }

  return nil
}

func main() {
  cmdArgs := parseArgs()

  setUpEnv(cmdArgs)

  // setup the cache, even though it isn't needed (to even nil pointer derefence errors in some places
	cache.LoadJSCache("", true)

  if err := createGraph(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }
}
