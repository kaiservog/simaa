package main

func FindArgument(arg string, args[] string) string {
  for index, element := range args {
    if element == arg {
      return args[index+1]
    }
  }
  return ""
}

func ContainArgument(arg string, args[] string) bool {
  for _, element := range args {
    if element == arg {
      return true
    }
  }
  return false
}
