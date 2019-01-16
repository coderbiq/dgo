package vo

// ValueObject 定义值对象模型的外观
type ValueObject interface {
	Equal(other ValueObject) bool
	String() string
	Empty() bool
}

// Identity 定义标识模型的外观
type Identity interface {
	ValueObject
}
