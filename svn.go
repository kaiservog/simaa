package main

import "fmt"
import "os/exec"
import "bytes"
import "strings"

func Checkout(svnUrl, folder string) error {
  fmt.Println("Conectando no SVN", svnUrl, "...")
  cmd := exec.Command("svn", "checkout", svnUrl, folder)
	var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
	if err != nil {
    return err
	} else {
    fmt.Println("Download da branch realizado em", folder)
    return nil
  }
}

func TestSubversion() error {
  cmd := exec.Command("svn", "--version")
  err := cmd.Run()
  return err
}

func IsSubversionPath(path string) bool {
  return strings.Contains(path[:4], "http")
}
