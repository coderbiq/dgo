package model

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
)

// IDGenerator 存储当前使用的标识生成器
// 主要用于为消息生成唯一标识，也可用于为聚合或实体生成标识。
var IDGenerator IdentityGenerator = NewSnowflakeIDGenerator(1)

// IdentityGenerator 定义标识生成器
type IdentityGenerator interface {
	StringID() StringID
	LongID() LongID
}

// StringID 字符串类型的标识模型
type StringID string

// Equal 判断提供的另一个标识是否与当前标识为同一个标识
func (id StringID) Equal(other ValueObject) bool {
	otherStringID, ok := other.(StringID)
	if !ok {
		return false
	}
	return string(id) == string(otherStringID)
}

func (id StringID) String() string {
	return string(id)
}

// Empty 返回当前标识是否为一个空值
func (id StringID) Empty() bool {
	return string(id) == ""
}

// LongID 长整数类型的标识模型
type LongID int64

// Equal 判断提供的另一个标识是否与当前标识为同一个标识
func (id LongID) Equal(other ValueObject) bool {
	otherLongID, ok := other.(LongID)
	if !ok {
		return false
	}
	return int64(id) == int64(otherLongID)
}

func (id LongID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

// Empty 返回当前标识是否为一个空值
func (id LongID) Empty() bool {
	return int64(id) == 0
}

// Int64 返回标识的原始值
func (id LongID) Int64() int64 {
	return int64(id)
}

// SnowflakeIDGenerator 使用 Twitter snowflake 算法的标识生成器
type SnowflakeIDGenerator struct {
	node *snowflake.Node
}

// NewSnowflakeIDGenerator 创建一个 snowflake 算法的标识生成器
func NewSnowflakeIDGenerator(n int64) *SnowflakeIDGenerator {
	node, err := snowflake.NewNode(n)
	if err != nil {
		panic(err)
	}
	return &SnowflakeIDGenerator{node: node}
}

// StringID 生成一个字符串标识
func (g SnowflakeIDGenerator) StringID() StringID {
	return StringID(g.LongID().String())
}

// LongID 生成一个 int64 标识
func (g SnowflakeIDGenerator) LongID() LongID {
	return LongID(g.node.Generate())
}
