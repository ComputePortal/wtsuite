package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"../../pkg/files"
	"../../pkg/parsers"
)

type CmdArgs struct {
	inputFile string
	xml       bool
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-?|-h|--help] [--xml] inputFile\n", os.Args[0])
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
		false,
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
			case "--xml":
				cmdArgs.xml = true
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
	if cmdArgs.xml {
		p, err := parsers.NewHTMLParser(cmdArgs.inputFile)
		if err != nil {
			printSyntaxErrorAndExit(err)
		}

		p.DumpTokens()
	} else {
		p, err := parsers.NewUIParser(cmdArgs.inputFile)
		if err != nil {
			printSyntaxErrorAndExit(err)
		}

		p.DumpTokens()
	}
}

func main() {
	cmdArgs := parseArgs()

	buildSyntaxTree(cmdArgs)
}
