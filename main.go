package main

import "fmt"
import "os"
import "errors"

func main() {
  command := os.Args[1]

  if command == "deploy" {
    err := deploy()
    if err != nil {
      fmt.Println(err)
    }
  } else if command == "hash" {
    path := FindArgument("-d", "--directory", os.Args, "")
    if path == "" {
      fmt.Println("-d --directory é obrigatório")
      return
    }
    file := GetDebName(path + "/caixa-extracash/main/target/")
    h, err := HashFile(path + "/caixa-extracash/main/target/" + file)
    if err != nil {
      fmt.Println("Erro ao gerar Hash", err)
    }
    fmt.Println("Hash da aplicação é:", h)
  } else if command == "teste" {
    err := testEnvironment()
    if err != nil {
      fmt.Println(err)
    }
  } else if command == "help" || command == "ajuda" {
    help()
  }
}

func help() {
  fmt.Println("Comandos:")
  fmt.Println("\t deploy -> Realiza deploy de uma branch svn ou diretório para um ATM.")
  fmt.Println("\t hash -> Exibe o hash de um pacote .deb do SIMAA.")
  fmt.Println("\t teste -> testa seu ambiente para verificar se você tem todos programas dependentes para execução do deploy.")
  fmt.Println("\t help -> Exibe esta mensagem de ajuda.")
  fmt.Println("\t ajuda -> Exibe esta mensagem de ajuda.")

  fmt.Println("Argumentos:")
  fmt.Println("\t -u --user -> Informa usuário do ATM, valor padrão 'root'")
  fmt.Println("\t -pw --password -> Informa senha do ATM, valor padrão 'caixa'")
  fmt.Println("\t -h --host -> Informa host(IP) do ATM, valor padrão ''")
  fmt.Println("\t -p --port -> Informa porta remota do ATM, valor padrão '22'")
  fmt.Println("\t -s --svn -> Informa endereço svn da branch da versão do SIMAA a ser gerado")
  fmt.Println("\t -d --directory -> Informa endereço local da branch da versão do SIMAA a ser gerado")

  fmt.Println("Exemplos:")

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
    branch := FindArgument("-s", "--svn", os.Args, "")
    directory := FindArgument("-d", "--directory", os.Args, "")
    host := FindArgument("-h", "--host", os.Args, "192.168.1.109")
    port := FindArgument("-p", "--port", os.Args, "22")
    user := FindArgument("-u", "--user", os.Args, "root")
    password := FindArgument("-pw", "--password", os.Args, "caixa")

    CleanWorkDirectory()

    workDir := GetWorkDirectory()
    if branch != "" {
      err := Checkout(branch, workDir)
      if err != nil {
        fmt.Println("Erro ao baixar branch pelo svn", branch)
        return err
      }
    } else if directory != "" {
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
    } else {
      return errors.New("Deve ser informado --svn ou --directory")
    }

    mvnPath :=  ConcatPath(workDir, "caixa-extracash")

    err := MavenCleanInstall(ConcatPath(mvnPath, "pom.xml"))
    if err != nil {
      fmt.Println("Erro ao realizar mvn clean install em ", mvnPath)
      return err
    }

    fmt.Println("Realizando deploy no ATM")

    server := SimaaServer{host, port, user, password}
    deb := GetDebName(mvnPath + "/main/target/")

    fmt.Println("Nome do .deb é", deb)
    CopyToServer(server, mvnPath + "/main/target/" + deb, "/tmp/" + deb)
    return nil
}
