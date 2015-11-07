package main

import (
	"flag"
	"fmt"
	"strings"

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
	generators.EnumSetupFlags()

	flag.Parse()
	if filePath == "" {
		fmt.Println("-source must be specified")
		return
	}

	if isDirectory(filePath) {
	} else {
		f := types.NewFile(filePath)
	}

	switch strings.ToLower(target) {
	case "enum":
		generators.DoGenerateEnum(f, nil)
	default:
		fmt.Sprintf("Unreconized target type '%s'\n", target)
	}
}
