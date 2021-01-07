package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

const (
  DEFAULT_OUTPUTFILE = "a.gv"
  DEFAULT_TYPE = "class"
)

var cmdParser *parsers.CLIParser = nil

type CmdArgs struct {
  graphType string // eg. class
  outputFile string // file needed by the graphviz utility 'dot' to create the visual, must be specified
  entryFile string // this is the entry point
  analyzedFiles map[string]string // (analyzing everything available from the package.json upwards would be too messy)

  verbosity int
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
    graphType: "class",
    outputFile: "",
    entryFile: "",
    analyzedFiles: make(map[string]string),
		verbosity:     0,
	}

  var positional []string = nil

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [options] --type <graph-type> --output <output-file> <input-files>\n", os.Args[0]),
    `Graph types:
  class                  Class inheritance, explicit interface implementation
  instance               Instance properties`,
    []parsers.CLIOption{
      parsers.NewCLIUniqueFile("o", "output" , "-o, --output    <output-file> Defaults to \"" + DEFAULT_OUTPUTFILE + "\" if not set", false, &(cmdArgs.outputFile)),
      parsers.NewCLIUniqueEnum("t", "type"   , "-t, --type      See below", []string{"class", "instance"}, &(cmdArgs.graphType)),
      parsers.NewCLICountFlag("v", ""        , "-v[v[v..]]      Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIRemaining(&positional),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

  switch cmdArgs.graphType {
  case "class":
    if len(positional) == 0 {
      printMessageAndExit("Error: graph type class requires at least one input file")
    }
  case "instance":
    if len(positional) == 0 {
      printMessageAndExit("Error: graph type instance requires at least one input file")
    }
  default:
    printMessageAndExit("Error: unrecognized --type " + cmdArgs.graphType)
  }

  orderedInputFiles := make([]string, 0)

  for _, arg := range positional {
    info, err := os.Stat(arg)
    if os.IsNotExist(err) {
      printMessageAndExit("Error: \"" + arg + "\" not found")
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
        printMessageAndExit("Error: " + err.Error())
      }
    } else {
      absPath, err := filepath.Abs(arg)
      if err != nil {
        printMessageAndExit("Error: " + err.Error())
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

  return files.ResolvePackages(cmdArgs.entryFile)
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
