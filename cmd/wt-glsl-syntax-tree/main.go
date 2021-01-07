package main

import (
  "fmt"
  "os"

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
	cmdArgs := CmdArgs{
		"",
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [-?|-h|--help] <input-file>\n", os.Args[0]),
    "",
    []parsers.CLIOption{},
    parsers.NewCLIFile("", "", "", true, &(cmdArgs.inputFile)),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
}

func buildSyntaxTree(cmdArgs CmdArgs) {
	p, err := parsers.NewGLSLParser(cmdArgs.inputFile)
	if err != nil {
		printSyntaxErrorAndExit(err)
	}

	p.DumpTokens()
}

func main() {
  cmdArgs := parseArgs()

  setUpEnv(cmdArgs)

  buildSyntaxTree(cmdArgs)
}
