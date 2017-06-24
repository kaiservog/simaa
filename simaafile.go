package main

import (
 "fmt"
 "os"
 "io"
 "io/ioutil"
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
    os.Mkdir(filePath, os.FileMode(0777))
  }
}

func GetLastFolderFromPath(branchPath string) string {
  folders := strings.Split(branchPath, "/")
  if len(folders) > 2 {
    return folders[len(folders) -1]
  } else {
    folders = strings.Split(branchPath, "\\")
    return folders[len(folders) -1]
  }
}

func CopyFile(source string, dest string) (err error) {
    sourcefile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourcefile.Close()

    destfile, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destfile.Close()

    _, err = io.Copy(destfile, sourcefile)
    if err == nil {
        _, err := os.Stat(source)
        if err != nil {
            err = os.Chmod(dest, 0777)
        }
    }
    SetWritable(dest)
    return
}

func MakeParentFolder(source string, dest string) (string, error) {
  dir := ConcatPath(dest, GetLastFolderFromPath(source))
  err := os.Mkdir(dir, os.FileMode(0777))
  SetWritable(dir)

  if err != nil {
    return "", err
  }

  dest = ConcatPath(dest, GetLastFolderFromPath(source))
  return dest, nil
}

func SetWritable(filepath string) error {
 	err := os.Chmod(filepath, 0222)
 	return err
}

func CopyDir(source string, dest string) (err error) {
    _, err = os.Stat(source)
    if err != nil {
        return err
    }

    err = os.MkdirAll(dest, 0777)
    if err != nil {
        return err
    }

    SetWritable(dest)

    directory, _ := os.Open(source)

    objects, err := directory.Readdir(-1)

    for _, obj := range objects {
        sourcefilepointer := source + "/" + obj.Name()
        destinationfilepointer := dest + "/" + obj.Name()

        if obj.IsDir() {
            err = CopyDir(sourcefilepointer, destinationfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            err = CopyFile(sourcefilepointer, destinationfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        }
    }
    return
}

func GetDebName(path string) string {
  files, _ := ioutil.ReadDir(path)
  for _, f := range files {
    fileExtension := f.Name()[len(f.Name())-3:]
    if strings.Contains(fileExtension, "deb") {
      return f.Name()
    }
  }
  return ""
}
