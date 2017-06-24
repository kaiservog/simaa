package main

import "fmt"
import "os/exec"
import "bytes"

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
  }
}

func TestSubversion() error {
  cmd := exec.Command("svn", "--version")
  err := cmd.Run()
  return err
}
