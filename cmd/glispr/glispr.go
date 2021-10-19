package main

import (
	"glispr/internal"
	"os"
)

func main() {
	parser := internal.NewParser(os.Stdin)
	internal.Eval(parser.Expression())
}
