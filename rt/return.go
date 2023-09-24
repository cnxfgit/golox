package rt

import (
	"golox/object"
)

type Return struct {
	message string
	Value   object.Object
}

func (r *Return) Error() string {
	return r.message
}
