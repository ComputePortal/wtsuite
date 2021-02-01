package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/computeportal/wtsuite/pkg/parsers"
)

var cmdParser *parsers.CLIParser = nil

type CmdArgs struct {
	inputFile string
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
	// default args
	cmdArgs := CmdArgs{
		"",
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [-?|-h|--help] <input-file>\n", os.Args[0]),
    "Note: this tool can only be used to analyze the attribute syntax-tree",
    []parsers.CLIOption{},
    parsers.NewCLIFile("", "", "", true, &(cmdArgs.inputFile)),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

	return cmdArgs
}

func buildSyntaxTree(cmdArgs CmdArgs) {
  path := cmdArgs.inputFile
  if !filepath.IsAbs(path) {
    panic("path should be absolute")
  }

  p, err := parsers.NewTemplateParser(path)
  if err != nil {
    printSyntaxErrorAndExit(err)
  }

  p.DumpTokens()
}

func main() {
	cmdArgs := parseArgs()

	buildSyntaxTree(cmdArgs)
}
