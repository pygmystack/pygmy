package color

import (
	"fmt"

	"github.com/mattn/go-colorable"
)

var colorableOutput = colorable.NewColorableStdout()

func Print(input interface{}) {
	fmt.Fprint(colorableOutput, input)
}
