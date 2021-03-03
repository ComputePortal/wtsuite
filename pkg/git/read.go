package git

import (
  "os"
  "path/filepath"
  //billy       "gopkg.in/src-d/go-billy.v4"
  //billyosfs   "gopkg.in/src-d/go-billy.v4/osfs"

  gitcore      "gopkg.in/src-d/go-git.v4"
)

func readWorktree(wt *gitcore.Worktree, srcDir string) error {
  if err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    if info.IsDir() {
      return nil
    }

    if _, err := wt.Add(path); err != nil {
      return err
    }

    return nil
  }); err != nil {
    return err
  }

  return nil
}
