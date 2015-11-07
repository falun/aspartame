package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/falun/aspartame/generators"
)

func main() {
	var (
		target   string
		filePath string
		output   string
	)

	flag.StringVar(&target, "type", "enum", "What kind of sweetening should we be doing")
	flag.StringVar(&filePath, "source", "./", "File, or directory, to operate on")
	flag.StringVar(&output, "output", "stdout", "Where should we produce output (file[:name]|stdout)")

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

	f := g.LocateFile(filePath)
	if f == nil {
		log.Fatal("Could not find a file to operate on.")
	}

	g.DoGenerate(f, nil)
}
