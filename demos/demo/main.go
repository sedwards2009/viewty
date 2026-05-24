package main

import (tt "github.com/sedwards2009/termtronic" )

func main() {
  app := tt.NewApplication()
  app.EnableLogging(true)

  app.Run()
}
