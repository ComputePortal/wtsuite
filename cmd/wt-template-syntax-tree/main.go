package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
)

type CmdArgs struct {
	inputFile string
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-?|-h|--help] inputFile\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  Note: this tool can only be used to analyze the attribute syntax-tree\n")
	os.Exit(1)
}

func messageAndExit(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
	printUsageAndExit()
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	// default args
	cmdArgs := CmdArgs{
		"",
	}

	positional := make([]string, 0)

	i := 1
	n := len(os.Args)

	for i < n {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-?", "-h", "--help":
				printUsageAndExit()
			default:
				messageAndExit("Error: unrecognized flag " + arg)
			}
		} else {
			positional = append(positional, arg)
		}

		i++
	}

	if len(positional) != 1 {
		messageAndExit("Error: expected 1 positional argument")
	}

	cmdArgs.inputFile = positional[0]

	if err := files.AssertFile(cmdArgs.inputFile); err != nil {
		messageAndExit("Error: input file '" + cmdArgs.inputFile + "' " + err.Error())
	}

	inputFile, err := filepath.Abs(cmdArgs.inputFile)
	if err != nil {
		messageAndExit("Error: input file '" + cmdArgs.inputFile + "' " + err.Error())
	} else {
		cmdArgs.inputFile = inputFile
	}

	return cmdArgs
}

func buildSyntaxTree(cmdArgs CmdArgs) {
  path := cmdArgs.inputFile
  if !filepath.IsAbs(path) {
    panic("path should be absolute")
  }

  source := files.NewDefaultUIFileSource()

  p, err := parsers.NewUIParser(source, path)
  if err != nil {
    printSyntaxErrorAndExit(err)
  }

  p.DumpTokens()
}

func main() {
	cmdArgs := parseArgs()

	buildSyntaxTree(cmdArgs)
}
