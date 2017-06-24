package main

import "fmt"
import "os"


func main() {
  command := os.Args[1]

  if command == "deploy" {
    branch := os.Args[2]
    server := os.Args[3]

    CleanWorkDirectory()

    workDir := GetWorkDirectory()
    err := Checkout(branchPath, workDir)
    if err != nil {
      fmt.Println("Erro ao baixar branch pelo svn", branchPath)
      fmt.Println(err)
      return
    }

    mvnPath := ConcatPath(workDir, GetLastFolderFromPath(branchPath))
    err = MavenCleanInstall(mvnPath)
    if err != nil {
      fmt.Println("Erro ao realizar mvn clean install em ", mvnPath)
      fmt.Println(err)
      return
    }

    fmt.Println("Realizando deploy " + file )
    //Copy(file, "C:\\Users\\cesar\\Documents\\dev\\go\\src\\simaa\\output")
  } else if command == "remote-deploy" {
    file := os.Args[2]
    server := SimaaServer{"192.168.172.2", "22", "root", "caixa"}
    CopyToServer(server, file, "/tmp/testando.txt")
  } else if command == "hash" {
    file := os.Args[2]
    fmt.Println(HashFile(file))
  } else if command == "checkout" {
    branchPath := os.Args[2]
    Checkout(branchPath, GetWorkDirectory())
  } else if command == "work" {
    CleanWorkDirectory()
    fmt.Println(GetWorkDirectory())
  } else if command == "teste" {
    err := TestMaven()
    if err != nil {
      fmt.Println("Erro ao executar o comando 'mvn'", err)
      return
    }

    err = TestSubversion()
    if err != nil {
      fmt.Println("Erro ao executar o comando 'svn'", err)
      return
    }

    fmt.Println("Parabéns, Aparentemente tudo esta OK")
  }
}
