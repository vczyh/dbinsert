package main

import (
	_ "embed"
	"github.com/vczyh/dbinsert/cmd"
)

//go:embed VERSION
var version string

func main() {
	cmd.Execute(version)
}
