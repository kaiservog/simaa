package main

import "fmt"
import "os"


func main() {
  command := os.Args[1]

  if command == "deploy" {
    branch := os.Args[2]
    //server := os.Args[3]

    CleanWorkDirectory()

    workDir := GetWorkDirectory()
    if IsSubversionPath(branch) {
      err := Checkout(branch, workDir)
      if err != nil {
        fmt.Println("Erro ao baixar branch pelo svn", branch)
        fmt.Println(err)
        return
      }
    } else {
      _, err := os.Stat(branch)
      if err != nil {
        panic(err)
      }

      fmt.Println("Copiando fonte para Área de Trabalho", branch, "->", workDir)
      workDir, err = MakeParentFolder(branch, workDir)

      if err != nil {
        fmt.Println(err)
        return
      }

      err = CopyDir(branch, workDir)
      if err != nil {
        fmt.Println("Erro ao copiar código fonte", branch)
        fmt.Println(err)
        return
      } else {
        fmt.Println("Código fonte copiado com sucesso")
      }
    }

    mvnPath :=  ConcatPath(workDir, "caixa-extracash")

    err := MavenCleanInstall(ConcatPath(mvnPath, "pom.xml"))
    if err != nil {
      fmt.Println("Erro ao realizar mvn clean install em ", mvnPath)
      fmt.Println(err)
      return
    }

    fmt.Println("Realizando deploy " + mvnPath)
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
