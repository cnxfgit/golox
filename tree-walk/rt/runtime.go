package rt

import (
	"fmt"
	"golox/tree-walk/token"
	"os"
	"strconv"
)

var HadError bool = false
var HadRuntimeError bool = false

func ErrorLine(line uint, message string) {
	report(line, "", message)
}

func ErrorToken(t token.Token, message string) {
	if t.Type == token.Eof {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func report(line uint, where string, message string) {
	_, _ = fmt.Fprintln(os.Stderr, "[line "+strconv.Itoa(int(line))+"] Error"+where+": "+message)
	HadError = true
}
