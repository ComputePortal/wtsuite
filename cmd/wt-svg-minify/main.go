package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"../../pkg/directives"
	"../../pkg/parsers"
	"../../pkg/tree"
	"../../pkg/tree/styles"
)

type CmdArgs struct {
	fname         string
	output        string // if empty -> write to stdout
	humanReadable bool
	absPrecision  int
	relPrecision  int
	genPrecision  int
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] svg-file\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `	
Options:
  -?, -h, --help         Show this message
  -o <file>              Output file instead of stdout
  --human                With newlines and unnecessary spaces, for debugging of minifier
	--abs-precision <int>  Precision of positions wrt. viewbox (default is 4)
	--rel-precision <int>  Precision of relative path motions wrt. viewbox (default is 6)
	--gen-precision <int>  Precision of general floating point number (default is 2)
`)

	os.Exit(1)
}

func printMessage(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
}

func printMessageAndExit(msg string, printUsage bool) {
	printMessage(msg)
	if printUsage {
		printUsageAndExit()
	} else {
		os.Exit(1)
	}
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
		fname:         "",
		output:        "",
		humanReadable: false,
		absPrecision:  4,
		relPrecision:  6,
		genPrecision:  2,
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
			case "-o":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.output != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.output = os.Args[i+1]
					i++
				}
			case "--human":
				if cmdArgs.humanReadable {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.humanReadable = true
				}
			case "--abs-precision":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					x64, err := strconv.ParseInt(os.Args[i+1], 10, 64)
					x := int(x64)
					if err != nil || x < 0 {
						printMessageAndExit("Error: bad integer argument after "+arg, true)
					}

					cmdArgs.absPrecision = x
					i++
				}
			case "--rel-precision":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					x64, err := strconv.ParseInt(os.Args[i+1], 10, 64)
					x := int(x64)
					if err != nil || x < 0 {
						printMessageAndExit("Error: bad integer argument after "+arg, true)
					}

					cmdArgs.relPrecision = x
					i++
				}
			case "--gen-precision":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					x64, err := strconv.ParseInt(os.Args[i+1], 10, 64)
					x := int(x64)
					if err != nil || x < 0 {
						printMessageAndExit("Error: bad integer argument after "+arg, true)
					}

					cmdArgs.genPrecision = x
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

	cmdArgs.fname = positional[0]
	fname, err := filepath.Abs(cmdArgs.fname)
	if err != nil {
		printMessageAndExit("Error: svg file '"+cmdArgs.fname+"' "+err.Error(), true)
	} else {
		cmdArgs.fname = fname
	}

	if cmdArgs.output != "" {
		cmdArgs.output, err = filepath.Abs(cmdArgs.output)
		if err != nil {
			printMessageAndExit("Error: output file '"+cmdArgs.fname+"' "+err.Error(), true)
		}
	}

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
	// always compress the numbers
	tree.COMPRESS_NUMBERS = true

	if !cmdArgs.humanReadable {
		tree.NL = ""
		tree.TAB = ""

		styles.LAST_SEMICOLON = ""
	}

	tree.ABS_PRECISION = cmdArgs.absPrecision
	tree.REL_PRECISION = cmdArgs.relPrecision
	tree.GEN_PRECISION = cmdArgs.genPrecision
}

func buildSVGFile(path string) (string, error) {
	p, err := parsers.NewHTMLParser(path)
	if err != nil {
		return "", err
	}

	rawTags, err := p.BuildTags()
	if err != nil {
		return "", err
	}

	root := tree.NewSVGRoot(p.NewContext(0, 1))
	node := directives.NewRootNode(root, directives.SVG)
	fileScope := directives.NewRootScope(false)

	for _, tag := range rawTags {
		if err := directives.BuildTag(fileScope, node, tag); err != nil {
			return "", err
		}
	}

	root.FoldDummy() // just to be sure that dummy tag isnt used

	tree.RegisterParents(root)

	// compression of svg child is done during write
	if err := root.Validate(); err != nil {
		return "", err
	}

	root.Minify()

	return root.Write("", tree.NL, tree.TAB), nil
}

func main() {
	cmdArgs := parseArgs()

	setUpEnv(cmdArgs)

	result, err := buildSVGFile(cmdArgs.fname)
	if err != nil {
		printMessageAndExit(err.Error(), false)
	}

	if cmdArgs.output == "" {
		fmt.Fprintf(os.Stdout, result)
	} else {
		if err := ioutil.WriteFile(cmdArgs.output, []byte(result), 0644); err != nil {
			printMessageAndExit("Error: "+err.Error(), false)
		}
	}
}
