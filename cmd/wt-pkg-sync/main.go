package main

// tag a commit as follows:
//  git tag -a v1.2.3 -m "release v1.2.3"
import (
  "fmt"
  "os"

  git    "gopkg.in/src-d/go-git.v4"
  gconf  "gopkg.in/src-d/go-git.v4/config"
  memory "gopkg.in/src-d/go-git.v4/storage/memory"
  memfs  "gopkg.in/src-d/go-billy.v4/memfs"
)

const TEST_URL = "https://github.com/computeportal/wtsuite.git"

func printMessageAndExit(msg string) {
  fmt.Fprintf(os.Stderr, msg)
  os.Exit(1)
}

func testClone() error {
  fs := memfs.New()

  storer := memory.NewStorage()

  _, err := git.Clone(storer, fs, &git.CloneOptions{
    URL: TEST_URL,
  })

  if err != nil {
    return err
  }

  rootFiles, err := fs.ReadDir("/")
  if err != nil {
    return err
  }

  for _, rootFile := range rootFiles {
    fmt.Println("found file: " + rootFile.Name())
  }

  return nil
}

func testListTags() error {
  storer := memory.NewStorage()

  fmt.Println("loading...")
  remote := git.NewRemote(storer, &gconf.RemoteConfig{
    URLs: []string{TEST_URL},
  })

  lst, err := remote.List(&git.ListOptions{})
  if err != nil {
    return err
  }
  fmt.Println("done loading")

  for _, ref := range lst {
    refName := ref.Name()
    fmt.Println("found ref: ", string(refName), refName.IsTag())
  }

  return nil
}

func main() {
  //if err := testClone(); err != nil {
    //printMessageAndExit(err.Error())
  //}

  if err := testListTags(); err != nil {
    printMessageAndExit(err.Error())
  }
}
