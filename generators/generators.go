package generators

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/falun/aspartame/types"
)

type Generator interface {
	SetupFlags()
	DoGenerate(*types.File, io.Writer)
}

// Each of the types of generation we support
var Generators = map[string]Generator{
	"enum": EnumGenerator,
}

func All() []Generator {
	var result = make([]Generator, 0)

	for _, v := range Generators {
		result = append(result, v)
	}

	return result
}

func Find(named string) (Generator, error) {
	g, ok := Generators[strings.ToLower(named)]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Could not find Generator named '%s'", named))
	} else {
		return g, nil
	}
}
