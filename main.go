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

    fmt.Println("Realizando deploy")

    server := SimaaServer{"192.168.1.109", "22", "chip", "chip"}
    deb := GetDebName(mvnPath + "/main/target/")

    fmt.Println("Nome do .deb é", deb)
    CopyToServer(server, mvnPath + "/main/target/" + deb, "/tmp/" + deb)
  } else if command == "remote-deploy" {
    server := SimaaServer{"192.168.1.109", "22", "chip", "chip"}
    caminho := "C:\\Users\\cesar\\Documents\\dev\\SIMAA\\biometria-saldo-saque-tev\\caixa-extracash\\main\\target"
    deb := GetDebName(caminho)
    fmt.Println("Nome do .deb é", deb)
    CopyToServer(server, caminho + "/" + deb, "/tmp/" + deb)
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
