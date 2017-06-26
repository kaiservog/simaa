package main

import "fmt"
import "os"


func main() {
  command := os.Args[1]

  if command == "deploy" {
    err := deploy()
    if err != nil {
      fmt.Println(err)
    }
  } else if command == "hash" {
    file := os.Args[2]
    h, err := HashFile(file)
    if err != nil {
      fmt.Println("Erro ao gerar Hash", err)
    }
    fmt.Println("Hash da aplicação é:", h)
  } else if command == "teste" {
    err := testEnvironment()
    if err != nil {
      fmt.Println(err)
    }
  }
}
func testEnvironment() error {
  err := TestMaven()

  if err != nil {
    fmt.Println("Erro ao executar o comando 'mvn'", err)
    return err
  }

  err = TestSubversion()
  if err != nil {
    fmt.Println("Erro ao executar o comando 'svn'", err)
    return err
  }

  fmt.Println("Parabéns, Aparentemente tudo esta OK")
  return nil
}

func deploy() error {
    branch := os.Args[2]

    CleanWorkDirectory()

    workDir := GetWorkDirectory()
    if IsSubversionPath(branch) {
      err := Checkout(branch, workDir)
      if err != nil {
        fmt.Println("Erro ao baixar branch pelo svn", branch)
        return err
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
        return err
      }

      err = CopyDir(branch, workDir)
      if err != nil {
        fmt.Println("Erro ao copiar código fonte", branch)
        return err
      } else {
        fmt.Println("Código fonte copiado com sucesso")
      }
    }

    mvnPath :=  ConcatPath(workDir, "caixa-extracash")

    err := MavenCleanInstall(ConcatPath(mvnPath, "pom.xml"))
    if err != nil {
      fmt.Println("Erro ao realizar mvn clean install em ", mvnPath)
      return err
    }

    fmt.Println("Realizando deploy no ATM")

    server := SimaaServer{"192.168.1.109", "22", "chip", "chip"}
    deb := GetDebName(mvnPath + "/main/target/")

    fmt.Println("Nome do .deb é", deb)
    CopyToServer(server, mvnPath + "/main/target/" + deb, "/tmp/" + deb)
    return nil
}
