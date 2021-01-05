package main

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "regexp"
  "strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tree/shaders"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

var VERBOSITY = 0

const (
  DEFAULT_OUTPUTFILE = "a.shader"
)

type CmdArgs struct {
  inputFile string
  outputFile string // defaults to a.shader in current dir

  target string
  compactOutput bool

  verbosity int
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s <input-file> [-o <output-file>] [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `
Options:
  -?, -h, --help            Show this message, other options are ignored
	-c, --compact             Compact output with minimal whitespace and short names
	-o, --output <file>       Defaults to "a.js" if not set
  -t, --target <shaderType> "vertex" or "fragment", defaults to "vertex"
	-v[v[v[v...]]]            Verbosity
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
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
		inputFile:     "",
		outputFile:    "",
    target:        "",
		compactOutput: false,
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
			case "-c", "--compact":
				if cmdArgs.compactOutput {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.compactOutput = true
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
			case "-t", "--target":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.target != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.target = os.Args[i+1]
					i++
				}
      }
		} else {
			positional = append(positional, arg)
		}

		i++
	}

	if len(positional) != 1 {
		printMessageAndExit("Error: expected 1 positional argument", true)
	}

	cmdArgs.inputFile = positional[0]

	if cmdArgs.outputFile == "" {
		cmdArgs.outputFile = DEFAULT_OUTPUTFILE
	}

	if err := files.AssertFile(cmdArgs.inputFile); err != nil {
		printMessageAndExit("Error: input file '"+cmdArgs.inputFile+"' "+err.Error(), false)
	}

	inputFile, err := filepath.Abs(cmdArgs.inputFile)
	if err != nil {
		printMessageAndExit("Error: input file '"+cmdArgs.inputFile+"' "+err.Error(), false)
	} else {
		cmdArgs.inputFile = inputFile
	}

	outputFile, err := filepath.Abs(cmdArgs.outputFile)
	if err != nil {
		printMessageAndExit("Error: output file '"+cmdArgs.outputFile+"' "+err.Error(), false)
	} else {
		cmdArgs.outputFile = outputFile
	}

  if cmdArgs.target != "" {
    if cmdArgs.target != "vertex" && cmdArgs.target != "fragment" {
      printMessageAndExit("Error: expected \"vertex\" or \"fragment\" target, got " + cmdArgs.target, true)
    }
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
