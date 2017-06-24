package main

import (
 "golang.org/x/crypto/ssh"
 "fmt"
 "os"
 "bytes"
 "net"
 "github.com/pkg/sftp"
 "log"
 "io"
)

type SimaaServer struct {
  Host string
  Port string
  User string
  Password string
}
func (server SimaaServer) GetHostWithPort() string {
  return server.Host + ":" + server.Port
}

func CopyToServer(server SimaaServer, originFilePath, serverFilePath string) {
  fmt.Println("Conectando com o servidor", server.Host)
  client, err := ConnectSSH(server)
  fmt.Println("Conexao com o servidor", server.Host, "realizada com sucesso")
  sftp, err := sftp.NewClient(client)
  fmt.Println("Conexao SFTP com o servidor", server.Host, "realizada com sucesso")
  if err != nil {
		log.Fatal(err)
  }
  defer sftp.Close()

  serverFile, err := sftp.Create(serverFilePath)
  if err != nil {
		log.Fatal(err)
  }

  localFile, err := os.Open(originFilePath)
  if err != nil {
		log.Fatal(err)
  }

  //bytesLocalFile := make([]byte, 5)
  //localFile.Read(bytesLocalFile)
  fmt.Println("Copiando arquivo")

  if _, err = io.Copy(serverFile, localFile); err != nil {
    fmt.Println("Erro ao copiar arquivo")
    log.Fatal(err)
  }

  //if _, err := serverFile.Write(bytesLocalFile); err != nil {
	//	log.Fatal(err)
  //}

  _, err = sftp.Lstat(serverFilePath)
	if err != nil {
		log.Fatal(err)
	}

  fmt.Println("Arquivo copiado com sucesso")
}

func ConnectSSH(server SimaaServer) (*ssh.Client, error) {
  clientConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{ssh.Password("caixa")},
    HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
  }

  client, err := ssh.Dial("tcp", server.GetHostWithPort(), clientConfig)
  if err != nil {
    return nil, err
  }

  return client, nil
}

func CommandOnServer(server SimaaServer)  {
  client, err := ConnectSSH(server)

  if err != nil {
    fmt.Println("Erro ao conectar no servidor ", server)
    return
  }

  session, err := client.NewSession()
  if err != nil {
    fmt.Print("Erro ao conectar em: " + server.Host + server.Port)
    panic(err.Error())
  }
  defer session.Close()

  var consoleOutput bytes.Buffer
  session.Stdout = &consoleOutput
  err = session.Run("ls -l /opt")

  if err != nil {
    fmt.Print("Erro ao executar comando em: " + server.Host + server.Port)
  } else {
    fmt.Println(consoleOutput.String())
    fmt.Println("Tudo ok")
  }
}
