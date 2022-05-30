package setuphelpers

import (
	"io"
	"monkey/evaluator"
	"monkey/object"
)

const MONKE = `🙈`

func LoadBuiltInMethods(env *object.Environment) {
	for key, value := range evaluator.BUILTIN {
		env.Set(key, value)
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "\n"+MONKE+" Error!:\n")
	for _, msg := range errors {
		io.WriteString(out, "> "+msg+"\n\n")
	}
}
