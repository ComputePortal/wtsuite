package main

import (
  "fmt"
  "net/http"
  "os"
  "strconv"
  "time"

  "github.com/computeportal/wtsuite/pkg/parsers"
)

const DEFAULT_PORT = 8080

var (
  cmdParser *parsers.CLIParser = nil
)

type CmdArgs struct {
  root string
  port int // 

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
		root:      "",
		port:      DEFAULT_PORT,
		verbosity: 0,
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s <root> [options]\n", os.Args[0]),
    "Test webserver for static site.",
    []parsers.CLIOption{
      parsers.NewCLIUniqueInt("p", "port"       , "-p, --port          Localhost port", &(cmdArgs.port)),
      parsers.NewCLICountFlag("v", ""               , "-v[v[v..]]             Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIDir("", "", "", true, &(cmdArgs.root)),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

  if cmdArgs.port <= 0 {
    printMessageAndExit("Error: invalid port value " + strconv.Itoa(cmdArgs.port))
  }

  return cmdArgs
}

func serve(cmdArgs CmdArgs) error {
  handle, err := NewRouter(cmdArgs.root)
  if err != nil {
    return err
  }

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(cmdArgs.port),
		Handler:        handle,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

  return server.ListenAndServe()
}

func main() {
  cmdArgs := parseArgs()

  if err := serve(cmdArgs); err != nil {
    printMessageAndExit(err.Error())
  }
}
