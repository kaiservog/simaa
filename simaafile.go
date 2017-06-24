package main

import (
 "fmt"
 "os"
 "io"
 "strings"
 "path/filepath"
)

func ConcatPath(path, file string) string {
  return path + string(filepath.Separator) + file
}

func Copy(source, target string) {
  in, err := os.Open(source)
  if err != nil {
    fmt.Println("Erro ao procurar arquivo", source)
    return
  }
  defer in.Close()

  fileName := FileNameFromPath(source)
  target = target + string(filepath.Separator) + fileName

  out, err := os.Create(target)
  if err != nil {
    fmt.Println("Erro ao criar arquivo de destino", target)
    return
  }

  defer func() {
    if err != nil {
      fmt.Println("Erro ao criar arquivo de destino", target)
    }
  }()

  if _, err = io.Copy(out, in); err != nil {
    fmt.Println("Erro ao copiar arquivo", source)
    return
  }
}

func FileNameFromPath(filePath string) string {
  if strings.Contains(filePath, "/") {
    f := strings.Split(filePath, "/")
    return f[len(f) -1]
  } else if strings.Contains(filePath, "\\"){
    f := strings.Split(filePath, "\\")
    return f[len(f) -1]
  } else {
    return filePath
  }
}

func CleanWorkDirectory() {
  directory := GetWorkDirectory()
  err := os.RemoveAll(directory)
  if err != nil {
    fmt.Println("erro", err)
  }
  GetWorkDirectory()
}

func GetWorkDirectory() string {
  _, err := os.Stat("C:\\")
  if err == nil {
    createDirectory("C:\\tmp")
    createDirectory("C:\\tmp\\simaa_work")
    return "C:\\tmp\\simaa_work"
  } else {
    createDirectory("/tmp")
    createDirectory("/tmp/simaa_work")
    return "/tmp/simaa_work"
   }
}

func createDirectory(filePath string) {
  _, err := os.Stat(filePath)

  if err != nil && os.IsNotExist(err) {
    os.Mkdir(filePath, 0777)
  }
}

func GetLastFolderFromPath(string branchPath) string {
  folders := strings.Split(branchPath, "/")
  if len(folders) > 2 {
    return folders[len(folders) -1]
  } else {
    folders = strings.Split(branchPath, "\\")
    return folders[len(folders) -1]
  }
}
