package model

// ValueObject model
type ValueObject interface {
	Equal(other ValueObject) bool
	String() string
	Empty() bool
}

// Identity model
type Identity interface {
	ValueObject
}
