package plugin

import (
	"github.com/Breeze0806/go-etl/element"
)

//RecordComparable 记录比较器
type RecordComparable interface {
	//比较记录left和right的差异，若列数不同，会返回错误
	Cmp(left, right element.Record) (Differ, error)
}

//DifferType 差异类型
type DifferType int8

//差异类型枚举
const (
	DifferTypeNone       DifferType = iota //两边无差异
	DifferTypeOnlyMaster                   //只有主数据库有
	DifferTypeOnlySlave                    //只有从数据库有
	DifferTypeValue                        //两边列值不同
)

var differTypeMap = map[DifferType]string{
	DifferTypeNone:       "none",
	DifferTypeOnlyMaster: "only master",
	DifferTypeOnlySlave:  "only slave",
	DifferTypeValue:      "value differ",
}

var differTypeChineseMap = map[DifferType]string{
	DifferTypeNone:       "无差异",
	DifferTypeOnlyMaster: "只有主数据库有",
	DifferTypeOnlySlave:  "只有从数据库有",
	DifferTypeValue:      "值不同",
}

func (d DifferType) String() string {
	if s, ok := differTypeMap[d]; ok {
		return s
	}
	return "unknown"
}

//Chinese 差异类型中文表达
func (d DifferType) Chinese() string {
	if s, ok := differTypeChineseMap[d]; ok {
		return s
	}
	return "未知差异"
}

//Differ 差异
type Differ struct {
	Type         DifferType     //标识差异类型
	MasterColumn element.Column //类型为DifferTypeOnlyMaster或DifferTypeValue时，主数据库列有值
	SlaveColumn  element.Column //类型为DifferTypeOnlySlave或DifferTypeValue时，从数据库列有值
}
