package main

// tag a commit as follows:
// * for current  commit: git tag -a v1.2.3 -m "release v1.2.3"
// * for an older commit: git tag -a v1.2.3 -m "release v1.2.3" f9c72e....
// and don't forget to share afterwards:
//  git push origin --tags
import (
  "bufio"
  "errors"
  "fmt"
  "io"
  "os"
  "path/filepath"
  "strings"
  
  "github.com/computeportal/wtsuite/pkg/files"
  "github.com/computeportal/wtsuite/pkg/parsers"

  git         "gopkg.in/src-d/go-git.v4"
  gitconfig   "gopkg.in/src-d/go-git.v4/config"
  gitplumbing "gopkg.in/src-d/go-git.v4/plumbing"
  gitmemory   "gopkg.in/src-d/go-git.v4/storage/memory"
  billy       "gopkg.in/src-d/go-billy.v4"
  billymemfs  "gopkg.in/src-d/go-billy.v4/memfs"
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

func parseTag(ref_ gitplumbing.ReferenceName) (*files.SemVer, error) {
  ref := string(ref_)

  parts := strings.Split(ref, "/")

  if parts[1] != "tags" {
    return nil, errors.New("unhandled tag format")
  }

  tag := parts[2]
  if strings.HasPrefix(tag, "v") {
    tag = tag[1:]
  }

  return files.ParseSemVer(tag)
}

func selectTag(url string, semVerRange *files.SemVerRange) (gitplumbing.ReferenceName, error) {
  var res gitplumbing.ReferenceName
  storer := gitmemory.NewStorage()

  remote := git.NewRemote(storer, &gitconfig.RemoteConfig{
    URLs: []string{url},
  })

  // do we need to fetch all the tags in order to be able to make a proper selection?
  /*if err := remote.Fetch(&git.FetchOptions{
    Tags: git.AllTags,
  }); err != nil {
    return res, err
  }*/

  lst, err := remote.List(&git.ListOptions{})
  if err != nil {
    return res, err
  }

  found := false
  var bestSemVer *files.SemVer = nil

  for _, ref_ := range lst {
    ref := ref_.Name()

    if !ref.IsTag() {
      continue
    }

    semVer, err := parseTag(ref)
    if err != nil {
      continue
    }

    if semVerRange.Contains(semVer) {
      if !found || (found && semVer.After(bestSemVer))  {
        bestSemVer = semVer
        res = ref
        found = true
      } 
    }
  }

  if !found {
    return res, errors.New("no valid tags found")
  }

  return res, nil
}

func writeFile(fs billy.Filesystem, src string, dst string) error {
  // TODO: only write the files that pass parser tests?
  fIn, err := fs.Open(src)
  if err != nil {
    return err
  }

  defer fIn.Close()

  fOut, err := os.Create(dst)
  if err != nil {
    return err
  }

  defer fOut.Close()

  wOut := bufio.NewWriter(fOut)

  if _, err := io.Copy(wOut, fIn); err != nil {
    return err
  }

  wOut.Flush()

  return nil
}

func writeDir(fs billy.Filesystem, dirSrc string, dirDst string) error {
  if err := os.MkdirAll(dirDst, 0755); err != nil {
    return err
  }

  files, err := fs.ReadDir(dirSrc)
  if err != nil {
    return err
  }

  for _, file := range files {
    src := fs.Join(dirSrc, file.Name())
    dst := filepath.Join(dirDst, file.Name())
    
    if file.IsDir() {
      if err := writeDir(fs, src, dst); err != nil {
        return err
      }
    } else {
      if err := writeFile(fs, src, dst); err != nil {
        return err
      }
    }
  }

  return nil
}

func writeMemFS(fs billy.Filesystem, dst string) error {
  return writeDir(fs, "/", dst)
}

func cloneTag(url string, ref gitplumbing.ReferenceName, dst string) error {
  fs := billymemfs.New()

  storer := gitmemory.NewStorage()

  repo, err := git.Clone(storer, fs, &git.CloneOptions{
    URL: url,
    ReferenceName: ref,
    SingleBranch: true,
    NoCheckout: true, // checkout follows further down
    RecurseSubmodules: git.NoRecurseSubmodules,
    Progress: nil,
  })
  if err != nil {
    return err
  }

  worktree, err := repo.Worktree()
  if err != nil {
    return err
  }

  if err := worktree.Checkout(&git.CheckoutOptions{
    Branch: ref,
  }); err != nil {
    return err
  }

  return writeMemFS(fs, dst)
}

func syncPackage(url_ string, svr *files.SemVerRange) error {
  url := "https://" + url_ + ".git"

  tagRef, err := selectTag(url, svr)
  if err != nil {
    return errors.New("Error: problem with " + url + " (" + err.Error() + ")")
  }

  dstBase := files.PkgInstallDst(url_)

  semVer, err := parseTag(tagRef)
  if err != nil {
    panic("should've been caught before")
  }

  if semVer.Write() == "0.0.0" {
    panic("something went wrong")
  }

  dst := filepath.Join(dstBase, semVer.Write())

  if files.IsFile(dst) {
    return errors.New("Error: destination " + dst + " is a file")
  }

  if FORCE || !files.IsDir(dst) {
    fmt.Println("Downloading " + url + "@" + semVer.Write())
    if err := cloneTag(url, tagRef, dst); err != nil {
      return err
    }
  }

  return err
}

func main() {
  parseArgs()

  pwd, err := os.Getwd()
  if err != nil {
    printMessageAndExit(err.Error())
  }

  if err := files.SyncPackages(pwd, syncPackage); err != nil {
    printMessageAndExit(err.Error())
  }
}
