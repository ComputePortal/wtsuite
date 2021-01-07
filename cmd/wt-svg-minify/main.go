package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
)

var cmdParser *parsers.CLIParser = nil

type CmdArgs struct {
	fname         string
	output        string // if empty -> write to stdout
	humanReadable bool
	absPrecision  int
	relPrecision  int
	genPrecision  int
}

func printMessage(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
}

func printMessageAndExit(msg string) {
	printMessage(msg)
  os.Exit(1)
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

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [options] <svg-file>\n", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFile("o", "output", "-o, --output <output-file>  Output file instead of stdout", false, &(cmdArgs.output)),
      parsers.NewCLIUniqueFlag("", "human", "Extra whitespace for readability", &(cmdArgs.humanReadable)),
      parsers.NewCLIUniqueInt("", "abs-precision", "--abs-precision <int>  Precision of positions wrt. viewbox (default is 4)", &(cmdArgs.absPrecision)),
      parsers.NewCLIUniqueInt("", "rel-precision", "--rel-precision <int>  Precision of relative path motions wrt. viewbox (default is 6)", &(cmdArgs.relPrecision)),
      parsers.NewCLIUniqueInt("", "gen-precision", "--gen-precision <int>  Precision of general floating point number (default is 2)", &(cmdArgs.genPrecision)),
    },
    parsers.NewCLIFile("", "", "", true, &(cmdArgs.fname)),
  )
    
  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
	// always compress the numbers
	tree.COMPRESS_NUMBERS = true

	if !cmdArgs.humanReadable {
		patterns.NL = ""
		patterns.TAB = ""
		patterns.LAST_SEMICOLON = ""
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
  // the source isn't really used, because the svg file doesnt contain import statements
	fileScope := directives.NewFileScope(false, directives.NewFileCache())

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

	return root.Write("", patterns.NL, patterns.TAB), nil
}

func main() {
	cmdArgs := parseArgs()

	setUpEnv(cmdArgs)

	result, err := buildSVGFile(cmdArgs.fname)
	if err != nil {
		printMessageAndExit(err.Error())
	}

	if cmdArgs.output == "" {
		fmt.Fprintf(os.Stdout, result)
	} else {
		if err := ioutil.WriteFile(cmdArgs.output, []byte(result), 0644); err != nil {
			printMessageAndExit("Error: "+err.Error())
		}
	}
}
