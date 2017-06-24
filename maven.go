package main

import "fmt"
import "os/exec"
import "bytes"

func MavenCleanInstall(projectPath, folder string) {
  fmt.Println("Tentando compilar com o Maven", projectPath)
  cmd := exec.Command("mvn", "clean", "install", "-f", projectPath)
	var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()
	if err != nil {
    fmt.Println(err)
    fmt.Printf("%q", out.String())
	} else {
    fmt.Println("Download da branch realizado em", folder)
  }
}

func TestMaven() error {
  cmd := exec.Command("mvn", "-version")
  err := cmd.Run()
  return err
}
