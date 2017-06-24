package main

import "fmt"
import "os/exec"
import "bytes"

func MavenCleanInstall(projectPath string) error {
  fmt.Println("Tentando compilar com o Maven", projectPath)
  cmd := exec.Command("mvn", "clean", "install", "-f", projectPath)
	var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
	if err != nil {
    fmt.Printf("%q", out.String())
    return err
	}
  return nil
}

func TestMaven() error {
  cmd := exec.Command("mvn", "-version")
  err := cmd.Run()
  return err
}
