package object

type ObjFunction struct {
	Name string
}

func NewFunction() *ObjFunction {
	return &ObjFunction{}
}
