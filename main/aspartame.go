package main

import (
	"flag"
	"os"

	"github.com/falun/aspartame/generators"
	"github.com/falun/aspartame/types"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		return
	}

	filePath := flag.Args()[0]

	f := types.NewFile(filePath)

	generators.GenerateEnum(f, os.Stdout)
}
