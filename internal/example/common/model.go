package common

import "github.com/coderbiq/dgo/model"

type (
	// Points 定义积分点数据模型
	Points uint

	// Account 定义积分账户聚合基础的外观
	Account interface {
		ID() model.LongID
		OwnerID() model.StringID
		Points() Points
		Deposit(points Points)
		Consume(points Points) error
	}
)

// BaseAccount 实现基本的积分账户模型
type BaseAccount struct {
	Identity      model.LongID
	OwnerIdentity model.StringID
	CurPoints     Points
}

// ID 返回积分账户标识
func (a BaseAccount) ID() model.LongID {
	return a.Identity
}

// OwnerID 返回积分账户所属的会员标识
func (a BaseAccount) OwnerID() model.StringID {
	return a.OwnerIdentity
}

// Points 返回积分账户中当前的可用积分
func (a BaseAccount) Points() Points {
	return a.CurPoints
}

// Inc 返回当前积分加上一个积分值后的新积分值
func (p Points) Inc(points Points) Points {
	return Points(uint(p) + uint(points))
}

// Dec 返回当前积分值减去一个积分值后的新积分值
func (p Points) Dec(points Points) Points {
	return Points(uint(p) - uint(points))
}

// GreaterThan 返回当前积分值是否比指定积分值大
func (p Points) GreaterThan(other Points) bool {
	return uint(p) > uint(other)
}

// Zero 返回当前积分值是否为零
func (p Points) Zero() bool {
	return p == 0
}
