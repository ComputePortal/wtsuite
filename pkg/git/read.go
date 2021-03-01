package git

import (
  billy       "gopkg.in/src-d/go-billy.v4"
  billyosfs   "gopkg.in/src-d/go-billy.v4/osfs"
)

func readWorktree(srcDir string) (billy.Filesystem, error) {
  return billyosfs.New(srcDir), nil
}
