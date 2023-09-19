package rt

import "golox/tree-walk/object"

type Return struct {
	message string
	Value   object.Object
}

func (r *Return) Error() string {
	return r.message
}
