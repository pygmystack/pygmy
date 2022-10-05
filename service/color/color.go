package color

import (
	"fmt"

	"github.com/mattn/go-colorable"
)

var colorableOutput = colorable.NewColorableStdout()

// Print will print text to an interface using a colour via go-colourable.
func Print(input interface{}) {
	fmt.Fprint(colorableOutput, input)
}
