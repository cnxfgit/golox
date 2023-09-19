package rt

import "golox/tree-walk/token"

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return e.Message
}
