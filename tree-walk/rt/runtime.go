package rt

import (
	"fmt"
	"os"
	"strconv"
)

var HadError bool = false
var HadRuntimeError bool = false

func ErrorLine(line uint, message string) {
	report(line, "", message)
}

func report(line uint, where string, message string) {
	_, _ = fmt.Fprintln(os.Stderr, "[line "+strconv.Itoa(int(line))+"] Error"+where+": "+message)
	HadError = true
}
