package compiler

import . "golox/vm/token"

type local struct {
	name       Token
	depth      int
	isCaptured bool
}
