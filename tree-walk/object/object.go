package object

import "strconv"

type Object interface {
	ToString() string
}

type Number float64

func (n Number) ToString() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type Boolean bool

func (b Boolean) ToString() string {
	return strconv.FormatBool(bool(b))
}

type String string

func (s String) ToString() string {
	return string(s)
}
