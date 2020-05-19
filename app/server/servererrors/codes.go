package servererrors

const (
	EMPTY    Code = "EMPTY"
	INVALID  Code = "INVALID"
	INTERNAL Code = "INTERNAL"
	ARGISNIL Code = "ARGUMENTISNIL"
)

type Code string
