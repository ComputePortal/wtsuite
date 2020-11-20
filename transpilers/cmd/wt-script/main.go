package main

import (
	"errors"
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
	"../../pkg/tokens/js/macros"
	"../../pkg/tokens/js/values"
	"../../pkg/tree/scripts"
)

const (
	DEFAULT_OUTPUTFILE = "a.js"
	DEFAULT_TARGET     = "nodejs"
)

var (
	VERBOSITY = 0
)

type CmdArgs struct {
	inputFile   string // entry script
	outputFile  string // defaults to a.js in current dir
	includeDirs []string
	target      string

	compactOutput bool
	forceBuild    bool // delete cache and start fresh
  executable    bool // create an executable

	verbosity int
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s <input-file> [-o <output-file>] [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `
Options:
  -?, -h, --help         Show this message, other options are ignored
	-c, --compact          Compact output with minimal whitespace and short names
	-f, --force            Force a complete project rebuild
	-I, --include <dir>    Append a search directory to HTMLPPPATH
	-o, --output <file>    Defaults to "a.js" if not set
	--target <js-target>   Defaults to "nodejs", other possibilities are "worker" or "browser"
  --executable           Create an executable with a node hashbang (target must be nodejs)
	-v[v[v[v...]]]         Verbosity
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
		includeDirs:   make([]string, 0),
		target:        "",
		compactOutput: false,
		forceBuild:    false,
    executable:    false,
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
			case "-f", "--force":
				if cmdArgs.forceBuild {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.forceBuild = true
				}
			case "-I", "--include":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					cmdArgs.includeDirs = append(cmdArgs.includeDirs, os.Args[i+1])
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
			case "--target":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.target != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.target = os.Args[i+1]
					i++
				}
			case "--executable":
        cmdArgs.executable = true
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

	if cmdArgs.target == "" {
		cmdArgs.target = DEFAULT_TARGET
	}

  if cmdArgs.executable && cmdArgs.target != DEFAULT_TARGET {
    printMessageAndExit("Error: --executable can only be used if target is nodejs", false)
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

	for _, includeDir := range cmdArgs.includeDirs {
		if err := files.AssertDir(includeDir); err != nil {
			printMessageAndExit("Error: include dir '"+includeDir+"' "+err.Error(), false)
		}
	}

	if !(cmdArgs.target == "nodejs" || cmdArgs.target == "worker" || cmdArgs.target == "browser") {
		printMessageAndExit("Error: invalid target", true)
	}

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
	files.JS_MODE = true

	if cmdArgs.compactOutput {
		js.NL = ""
		js.TAB = ""
		js.COMPACT_NAMING = true
		macros.COMPACT = true
	}

	js.TARGET = cmdArgs.target
	directives.ForceNewViewFileScriptRegistration()

	VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

	files.AppendIncludeDirs(cmdArgs.includeDirs)
}

func buildProject(cmdArgs CmdArgs) error {
	cache.LoadJSCache(cmdArgs.outputFile, cmdArgs.forceBuild)

	if cache.RequiresUpdate(cmdArgs.inputFile) {
		entryScript, err := scripts.NewInitFileScript(cmdArgs.inputFile, "")
		if err != nil {
			return err
		}

		bundle := scripts.NewFileBundle(map[string]string{})

		bundle.Append(entryScript)

		if err := bundle.Finalize(); err != nil {
			return err
		}

		content, err := bundle.Write()
		if err != nil {
			return err
		}

    if cmdArgs.executable {
      if err := ioutil.WriteFile(cmdArgs.outputFile, []byte("#!/usr/bin/env node\n"+content), 0755); err != nil {
        return errors.New("Error: " + err.Error())
      }
    } else {
      if err := ioutil.WriteFile(cmdArgs.outputFile, []byte(content), 0644); err != nil {
        return errors.New("Error: " + err.Error())
      }
    }

		cache.SaveCache(cmdArgs.outputFile)
	}

	return nil
}

func main() {
	cmdArgs := parseArgs()

	setUpEnv(cmdArgs)

	if err := buildProject(cmdArgs); err != nil {
		printSyntaxErrorAndExit(err)
	}
}
