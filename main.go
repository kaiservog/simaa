package main

import "fmt"
import "os"
import "errors"
import "time"
import "strings"

func main() {
  var command string
  if len(os.Args) > 1 {
    command = os.Args[1]
  }

  if command == "deploy" {
    err := deploy()
    if err != nil {
      fmt.Println(err)
    }
  } else if command == "hash" {
    hash, err := geraHashCommand()
    if(err != nil) {
      fmt.Println(err)
      return
    }

    fmt.Println("Hash da aplicação é:", hash)
  } else if command == "teste" {
    err := testEnvironment()
    if err != nil {
      fmt.Println(err)
    }
  } else if command == "deploy-cert" {
    copyCertToServer()
  } else if command == "restore-db" {
    restoreDbInServer()
  } else if command == "help" || command == "ajuda" {
    help()
  } else {
    fmt.Println("Comando não encontrado... exibindo ajuda!")
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

  err = Test7z()
  if err != nil {
    fmt.Println("Erro ao executar o comando '7z'", err)
    return err
  }

  fmt.Println("Parabéns, Aparentemente tudo esta OK")
  return nil
}

func deployDeb(directory string) {
  fmt.Println("Realizando deploy no ATM")

  host := FindArgument("-h", "--host", os.Args, "192.168.1.109")
  port := FindArgument("-p", "--port", os.Args, "22")
  user := FindArgument("-u", "--user", os.Args, "root")
  password := FindArgument("-pw", "--password", os.Args, "caixa")

  server := SimaaServer{host, port, user, password}
  deb := GetDebName(directory + "/caixa-extracash/main/target")

  fmt.Println("Nome do .deb é", deb)
  CopyToServer(server, directory + "/caixa-extracash/main/target" + deb,
    "/tmp/" + deb)

  fmt.Println("Aguardando 1 minuto para reboot do ATM")
  time.Sleep(time.Second)
}

func deploy() error {
    branch := FindArgument("-s", "--svn", os.Args, "")
    directory := FindArgument("-d", "--directory", os.Args, "")

    CleanWorkDirectory()
    workDir := GetWorkDirectory()

    if branch != "" {
      err := Checkout(branch, workDir)
      if err != nil {
        fmt.Println("Erro ao baixar branch pelo svn", branch)
        return err
      }
      workDir, err = MakeParentFolder(branch, workDir)

    } else if directory != "" {
      _, err := os.Stat(directory)
      if err != nil {
        panic(err)
      }

      fmt.Println("Copiando fonte para Área de Trabalho", directory, "->", workDir)
      workDir, err = MakeParentFolder(directory, workDir)

      if err != nil {
        fmt.Println(err)
        return err
      }

      err = CopyDir(directory, workDir)
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

    deployDeb(workDir)

    return nil
}

func copyCertToServer() {
  host := FindArgument("-h", "--host", os.Args, "192.168.1.109")
  port := FindArgument("-p", "--port", os.Args, "22")
  user := FindArgument("-u", "--user", os.Args, "root")
  password := FindArgument("-pw", "--password", os.Args, "caixa")

  server := SimaaServer{host, port, user, password}
  CopyToServer(server, "C:/opt/foton/multicanal/trustedcert-keystore.jks", "/opt/foton/multicanal/trustedcert-keystore.jks")
}

func restoreDbInServer() {
  host := FindArgument("-h", "--host", os.Args, "192.168.1.109")
  port := FindArgument("-p", "--port", os.Args, "22")
  user := FindArgument("-u", "--user", os.Args, "root")
  password := FindArgument("-pw", "--password", os.Args, "caixa")

  server := SimaaServer{host, port, user, password}
  CommandOnServer(server, "rm -rf /opt/foton/db/EXTRACASH.GDB")
  CommandOnServer(server, "cp -f /opt/foton/db_backup/*.GDB /opt/foton/db/EXTRACASH.GDB")
}

func GeraHashFromOpt(optPath string) (string, error) {
  path := optPath + "/foton/multicanal/lib/"
  var files []string
  files = append(files,
    path + "acessibilidade.jar",
    path + "bloqueiocartao.jar",
    path + "capitalizacao.jar",
    path + "cartaocredito.jar",
    path + "cartaomultiplo.jar",
    path + "common.jar",
    path + "consulta.jar",
    path + "contraordemprovisoria.jar",
    path + "contrata.jar",
    path + "creditodiretocaixa.jar",
    path + "depositario.jar",
    path + "depositavalores.jar",
    path + "desbloqueio.jar",
    path + "doafomezero.jar",
    path + "emitedocumento.jar",
    path + "emitesegundaviaderecibo.jar",
    path + "emv.jar",
    path + "extrato.jar",
    path + "images.jar",
    path + "infra.jar",
    path + "investimento.jar",
    path + "libfdk.jar",
    path + "main.jar",
    path + "mantemcheque.jar",
    path + "manutencao.jar",
    path + "pagadocumento.jar",
    path + "penhor.jar",
    path + "poupanca.jar",
    path + "saldo.jar",
    path + "saque.jar",
    path + "servico.jar",
    path + "servicosms.jar",
    path + "transferencia.jar",
    path + "viradadata.jar")

    return HashFiles(files)
}

func geraHashCommand() (string, error) {
  path := FindArgument("-d", "--directory", os.Args, "C:/opt")
  if path == "" {
    return "", errors.New("-d --directory é obrigatório")
  }

  if strings.Contains(path, ".deb") {
    fmt.Println("Gerando hash de um arquivo 'deb'")
    CleanWorkDirectory()
    workDir := GetWorkDirectory()
    CreateDirectory(workDir + "/hash")
    err := ExtractDebFiles(path, ConcatPath(workDir, "hash"))
    if err != nil {
      fmt.Println("Erro ao extrair", err)
    }

    path = workDir + "/hash/opt"
  } else {
    fmt.Println("Gerando hash de pasta opt", path)
  }

  hash, err := GeraHashFromOpt(path)
  if err != nil {
    fmt.Println("Erro ao gerar Hash em", path)
    fmt.Println(err)
  }

  return hash, nil
}
