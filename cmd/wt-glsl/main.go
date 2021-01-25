package main

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tree/shaders"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

const (
  DEFAULT_OUTPUTFILE = "a.shader"
)

var VERBOSITY = 0

var cmdParser *parsers.CLIParser = nil

type CmdArgs struct {
  inputFile string
  outputFile string // defaults to a.shader in current dir

  target string
  compactOutput bool

  verbosity int
}


func printMessageAndExit(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
  os.Exit(1)
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
		inputFile:     "",
		outputFile:    DEFAULT_OUTPUTFILE,
    target:        "vertex",
		compactOutput: false,
		verbosity:     0,
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s <input-file> [-o <output-file>] [options]", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFile("o", "output" , "-o, --output    <output-file> Defaults to \"" + DEFAULT_OUTPUTFILE + "\" if not set", false, &(cmdArgs.outputFile)),
      parsers.NewCLIUniqueFlag("c", "compact", "-c, --compact   Compact output with minimal whitespace and short names", &(cmdArgs.compactOutput)),
      parsers.NewCLIUniqueEnum("t", "target" , "-t, --target    \"vertex\" or \"fragment\", defaults to \"vertex\"", []string{"vertex", "fragment"}, &(cmdArgs.target)),
      parsers.NewCLICountFlag("v", ""        , "-v[v[v..]]      Verbosity", &(cmdArgs.verbosity)),
      parsers.NewCLIUniqueFlag("l", "latest" , "-l, --latest    Ignore max semver, use latest tagged versions of dependencies", &(files.LATEST)),
    },
    parsers.NewCLIFile("", "", "", true, &(cmdArgs.inputFile)),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) error {
  if cmdArgs.compactOutput {
    patterns.NL = ""
    patterns.TAB = ""
    patterns.COMPACT_NAMING = true
  }

  if cmdArgs.target != "" {
    glsl.TARGET = cmdArgs.target
  }

	VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	shaders.VERBOSITY = cmdArgs.verbosity

  return files.ResolvePackages(cmdArgs.inputFile)
}

func buildShader(cmdArgs CmdArgs) error {
  // dont bother caching, because shaders are expected to be relatively small
  entryShader, err := shaders.NewInitShaderFile(cmdArgs.inputFile)
  if err != nil {
    return err
  }

  bundle := shaders.NewShaderBundle()

  bundle.Append(entryShader)

  if err := bundle.Finalize(); err != nil {
    return err
  }

  content, err := bundle.Write(patterns.NL, patterns.TAB)
  if err != nil {
    return err
  }

  if err := ioutil.WriteFile(cmdArgs.outputFile, []byte(content), 0644); err != nil {
    return errors.New("Error: " + err.Error())
  }

  return nil
}

func main() {
  cmdArgs := parseArgs()
  
  if err := setUpEnv(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }

  if err := buildShader(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }
}
