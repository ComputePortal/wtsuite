package main

// tag a commit as follows:
// * for current  commit: git tag -a v1.2.3 -m "release v1.2.3"
// * for an older commit: git tag -a v1.2.3 -m "release v1.2.3" f9c72e....
// and don't forget to share afterwards:
//  git push origin --tags
import (
  "fmt"
  "os"
  
  "github.com/computeportal/wtsuite/pkg/files"
  "github.com/computeportal/wtsuite/pkg/git"
  "github.com/computeportal/wtsuite/pkg/parsers"
)

var (
  cmdParser *parsers.CLIParser = nil
  FORCE = false
)

func printMessageAndExit(msg string) {
  fmt.Fprintf(os.Stderr, "%s\n", msg)
  os.Exit(1)
}

func parseArgs() {
  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [options]", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFlag("f", "force", "-f, --force  Force (re)download of all dependencies", &FORCE),
      parsers.NewCLIUniqueFlag("l", "latest", "-l, --latest   Download latest tag, ignore min/max semver in package.json files", &(files.LATEST)),
    },
    nil,
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }
}

func main() {
  parseArgs()

  pwd, err := os.Getwd()
  if err != nil {
    printMessageAndExit(err.Error())
  }

  if err := files.SyncPackages(pwd, git.FetchPublicOrPrivate); err != nil {
    printMessageAndExit(err.Error())
  }
}
