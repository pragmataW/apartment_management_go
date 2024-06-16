package dto

type FlatAlreadyExists struct{
	Message string
}

func (e FlatAlreadyExists) Error() string {
	return e.Message
}

type ThereIsNoFlat struct{
	Message string
}

func (e ThereIsNoFlat) Error() string {
	return e.Message
}

type ThereIsNoDues struct{
	Message string
}

func (e ThereIsNoDues) Error() string {
	return e.Message
}