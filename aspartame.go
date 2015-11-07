package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/falun/aspartame/generators"
	"github.com/falun/aspartame/types"
)

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func main() {
	var target string
	var filePath string

	flag.StringVar(&target, "target", "enum", "What kind of sweetening should we be doing")
	flag.StringVar(&filePath, "source", "./", "File, or directory, to operate on")

	for _, e := range generators.All() {
		e.SetupFlags()
	}

	flag.Parse()
	if filePath == "" {
		fmt.Println("-source must be specified")
		return
	}

	g, gErr := generators.Find(target)
	if gErr != nil {
		log.Fatal(gErr)
	}

	if isDirectory(filePath) {
		fmt.Println("Attempting to operate on directory:", filePath)
	} else {
		f := types.NewFile(filePath)
		g.DoGenerate(f, nil)
	}
}
