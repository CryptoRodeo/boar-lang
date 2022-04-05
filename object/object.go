package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

/**
Dev notes:
- every value we encounter and evaluate will be represented using an Object interace

**/
