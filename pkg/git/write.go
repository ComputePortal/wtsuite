package git

import (
  "bufio"
  "io"
  "os"
  "path/filepath"

  billy       "gopkg.in/src-d/go-billy.v4"
)

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

func writeWorktree(fs billy.Filesystem, dst string) error {
  return writeDir(fs, "/", dst)
}
