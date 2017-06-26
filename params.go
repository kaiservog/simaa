package main

func FindArgument(argA string, argB string, args[] string, def string) string {
  for index, element := range args {
    if element == argA || element == argB {
      return args[index+1]
    }
    }
  return def
}

func ContainArgument(arg string, args[] string) bool {
  for _, element := range args {
    if element == arg {
      return true
    }
  }
  return false
}
