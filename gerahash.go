package main

import (
 "io"
 "os"
 "crypto/md5"
 "encoding/hex"
 "os/exec"
 "fmt"
)

func HashString(data string) (string, error) {
  hasher := md5.New()
  io.WriteString(hasher, data)
  return hex.EncodeToString(hasher.Sum(nil)), nil
}

func HashFile(filePath string) (string, error) {
  var result []byte
  fmt.Println("Gerando hash de do arquivo", filePath)

  file, err := os.Open(filePath)
  if err != nil {
    return "", err
  }
  defer file.Close()

  hasher := md5.New()
  if _, err := io.Copy(hasher, file); err != nil {
    return "", err
  }

  return hex.EncodeToString(hasher.Sum(result)), nil
}

func HashFiles(filePaths []string) (string, error) {
  hashs := ""

  for _, element := range filePaths {
    newHash, err := HashFile(element)
    if err != nil {
      return "", err
    } else {
      hashs = hashs + newHash
    }
  }

  return HashString(hashs)
}

func Test7z() error {
  cmd := exec.Command("7z")
  err := cmd.Run()
  return err
}
