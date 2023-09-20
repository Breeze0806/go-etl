// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbms

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/pingcap/errors"
)

const (
	minDuration time.Duration = -1 << 63
	maxDuration time.Duration = 1<<63 - 1
)

type splitRangeFetcher interface {
	fetchMax(ctx context.Context, splitTable database.Table) (element.Column, error)
	fetchMin(ctx context.Context, splitTable database.Table) (element.Column, error)
}

// SplitConfig 切分配置
type SplitConfig struct {
	Key string `json:"key"` //切分键
	//day（日）,min（分钟）,s（秒）,ms（毫秒）,us（微秒）,ns（纳秒）
	TimeAccuracy string     `json:"timeAccuracy"` //切分时间精度（默认为day）
	Range        SplitRange `json:"range"`        //切分范围
}

func (s *SplitConfig) fetchMin(ctx context.Context,
	splitTable database.Table) (c element.Column, err error) {
	if err = s.build(splitTable); err != nil {
		return
	}
	return s.Range.leftColumn(s.Key)
}

func (s *SplitConfig) fetchMax(ctx context.Context,
	splitTable database.Table) (c element.Column, err error) {
	if err = s.build(splitTable); err != nil {
		return
	}
	return s.Range.rightColumn(s.Key)
}

func (s *SplitConfig) build(splitTable database.Table) (err error) {
	if err = s.checkType(splitTable); err != nil {
		return err
	}
	return s.setLayout()
}

func (s *SplitConfig) setLayout() error {
	if s.Range.Type == string(element.TypeTime) {
		tl := &timeLayout{}
		tl.getLayout(s.TimeAccuracy)
		s.Range.Layout = tl.layout
		if s.Range.Layout == "" {
			return fmt.Errorf(
				"timeAccuracy(%v) should not be empty or valid when set range as non-empty",
				s.TimeAccuracy)
		}
	}
	return nil
}

func (s SplitConfig) checkType(splitTable database.Table) (err error) {
	if typ, ok := splitTable.Fields()[0].Type().(database.ValuerGoType); ok {
		switch typ.GoType() {
		case database.GoTypeInt64:
			switch element.ColumnType(s.Range.Type) {
			case element.TypeBigInt:
			default:
				return fmt.Errorf("checkType %v is not %v but %v",
					s.Key, database.GoTypeInt64, s.Range.Type)
			}
		case database.GoTypeString:
			switch element.ColumnType(s.Range.Type) {
			case element.TypeBigInt, element.TypeString:
			default:
				return fmt.Errorf("checkType %v is not %v but %v",
					s.Key, database.GoTypeString, s.Range.Type)
			}
		case database.GoTypeTime:
			switch element.ColumnType(s.Range.Type) {
			case element.TypeTime:
			default:
				return fmt.Errorf("checkType %v is not %v but %v",
					s.Key, database.GoTypeTime, s.Range.Type)
			}
		}
	}
	return nil
}

// SplitRange 切分范围配置
type SplitRange struct {
	Type   string `json:"type"`   //类型 bigint, string, time
	Layout string `json:"layout"` //时间格式
	Left   string `json:"left"`   //开始点
	Right  string `json:"right"`  //结束点
	where  string
}

func (s SplitRange) leftColumn(key string) (element.Column, error) {
	return s.fetchColumn(key, s.Left)
}

func (s SplitRange) rightColumn(key string) (element.Column, error) {
	return s.fetchColumn(key, s.Right)
}

func (s SplitRange) fetchColumn(key string, value string) (element.Column, error) {
	switch element.ColumnType(s.Type) {
	case element.TypeBigInt:
		bi, ok := new(big.Int).SetString(value, 10)
		if !ok {
			return nil, errors.Errorf("column(%v) value is not %v", key, element.TypeBigInt)
		}
		return element.NewDefaultColumn(element.NewBigIntColumnValue(bi), key, 0), nil
	case element.TypeString:
		return element.NewDefaultColumn(element.NewStringColumnValue(value), key, 0), nil
	case element.TypeTime:
		t, err := time.Parse(s.Layout, value)
		if err != nil {
			return nil, errors.Wrapf(err, "column(%v) value is not valid time", key)
		}
		return element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(t,
			element.NewStringTimeDecoder(s.Layout)), key, 0), nil
	}
	return nil, errors.Errorf("column(%v) type(%v) does not support", key, s.Type)
}

func split(min, max element.Column, num int,
	timeAccuracy string, splitField database.Field) (ranges []SplitRange, err error) {
	if num < 1 {
		err = errors.Errorf("splitNumber(%d) can not less than 1.", num)
		return
	}

	if min == nil || max == nil {
		err = errors.New("split min or max can not be nil")
		return
	}

	var left, right *big.Int
	var c convertor

	c, err = newConvertor(min, timeAccuracy)
	if err != nil {
		return
	}

	left, err = c.toBigInt(min)
	if err != nil {
		return
	}

	right, err = c.toBigInt(max)
	if err != nil {
		return
	}

	results := doSplit(left, right, num)

	for i := 0; i < len(results)-1; i++ {
		format := "%s >= %s and %s < %s"
		if i == len(results)-2 {
			format = "%s >= %s and %s <= %s"
		}

		typ, layout := c.splitConfig()
		ran := SplitRange{
			Type:   typ,
			Layout: layout,
			Left:   c.fromBigInt(results[i]),
			Right:  c.fromBigInt(results[i+1]),
			where: fmt.Sprintf(format, splitField.Quoted(), splitField.BindVar(1),
				splitField.Quoted(), splitField.BindVar(2)),
		}
		ranges = append(ranges, ran)
	}
	return
}

func newConvertor(min element.Column, timeAccuracy string) (convertor, error) {
	switch data := min.(*element.DefaultColumn).ColumnValue.(type) {
	case *element.BigIntColumnValue:
		return &bigIntConvertor{}, nil
	case *element.StringColumnValue:
		return &stringConvertor{}, nil
	case *element.TimeColumnValue:
		layout := &timeLayout{layout: data.Layout()}
		layout.getLayout(timeAccuracy)
		t, _ := data.AsTime()
		return &timeConvertor{layout: layout, min: t}, nil
	}
	return nil, errors.Errorf("split key can not be %v", min.Type())
}

type convertor interface {
	splitConfig() (typ string, layout string)
	fromBigInt(bi *big.Int) (v string)
	toBigInt(c element.Column) (bi *big.Int, err error)
}

type bigIntConvertor struct{}

func (b *bigIntConvertor) splitConfig() (typ string, layout string) {
	return element.TypeBigInt.String(), ""
}

func (b *bigIntConvertor) fromBigInt(bi *big.Int) (v string) {
	return bi.String()
}

func (b *bigIntConvertor) toBigInt(c element.Column) (bi *big.Int, err error) {
	var v element.BigIntNumber
	if v, err = c.AsBigInt(); err != nil {
		err = errors.Wrap(err, "AsBigInt fail")
		return
	}
	bi = v.AsBigInt()
	return
}

type stringConvertor struct{}

func (s *stringConvertor) splitConfig() (typ string, layout string) {
	return element.TypeString.String(), ""
}

func (s *stringConvertor) fromBigInt(bi *big.Int) (v string) {
	return bigint2String(bi, 128)
}

func (s *stringConvertor) toBigInt(c element.Column) (bi *big.Int, err error) {
	var v string
	if v, err = c.AsString(); err != nil {
		err = errors.Wrap(err, "AsString fail")
		return
	}
	return string2Bigint(v, 128)
}

func string2Bigint(s string, radix int64) (res *big.Int, err error) {
	res = big.NewInt(0)
	radixBigInt := big.NewInt(radix)
	for _, r := range s {
		if r < 0x0000 || r >= 0x0080 {
			return nil, errors.Errorf("split only can support ASCII. string[%s] is not ASCII string", s)
		}
		res = new(big.Int).Add(big.NewInt(int64(r)), new(big.Int).Mul(res, radixBigInt))
	}
	return res, nil
}

func bigint2String(res *big.Int, radix int64) string {
	var data []byte
	radixBigInt := big.NewInt(radix)
	zero := big.NewInt(0)
	for quotient := new(big.Int).Set(res); quotient.Cmp(zero) > 0; quotient = new(big.Int).Div(quotient, radixBigInt) {
		remainder := new(big.Int).Mod(quotient, radixBigInt)
		data = append(data, byte(remainder.Int64()))
	}
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-1-i] = data[len(data)-1-i], data[i]
	}
	return string(data)
}

type timeLayout struct {
	layout string
}

func (t *timeLayout) unit() time.Duration {
	switch len(t.layout) {
	case 10:
		return 24 * time.Hour
	case 16:
		return time.Minute
	case 19:
		return time.Second
	case 23:
		return time.Millisecond
	case 26:
		return time.Microsecond
	case 29:
		return time.Nanosecond
	default:
		return 0
	}
}

func (t *timeLayout) getLayout(timeAccuracy string) {
	switch timeAccuracy {
	case "day":
		t.layout = element.DefaultTimeFormat[:10]
	case "min":
		t.layout = element.DefaultTimeFormat[:16]
	case "s":
		t.layout = element.DefaultTimeFormat[:19]
	case "ms":
		t.layout = element.DefaultTimeFormat[:23]
	case "us":
		t.layout = element.DefaultTimeFormat[:26]
	case "ns":
		t.layout = element.DefaultTimeFormat[:29]
	}
}

type timeConvertor struct {
	layout *timeLayout
	min    time.Time
}

func (t *timeConvertor) splitConfig() (typ string, layout string) {
	return element.TypeTime.String(), t.layout.layout
}

func (t *timeConvertor) fromBigInt(bi *big.Int) (v string) {
	ti := t.min.Add(time.Duration(bi.Int64()) * t.layout.unit())
	return ti.Format(t.layout.layout)
}

func (t *timeConvertor) toBigInt(c element.Column) (bi *big.Int, err error) {
	var v time.Time
	if t.layout.unit() == 0 {
		return nil, errors.Errorf("time layout(%v) is not valid", t.layout.layout)
	}
	if v, err = c.AsTime(); err != nil {
		err = errors.Wrap(err, "AsTime fail")
		return
	}
	d := v.Sub(t.min)
	if d == minDuration || d == maxDuration {
		err = errors.Errorf("the grap (%v - %v) is too large", v, t.min)
		return
	}

	return big.NewInt(int64(d / t.layout.unit())), nil
}

func doSplit(left *big.Int, right *big.Int, num int) (results []*big.Int) {
	if left.Cmp(right) == 0 {
		results = []*big.Int{left, right}
		return
	}

	if left.Cmp(right) > 0 {
		left, right = right, left
	}

	gap := new(big.Int).Sub(right, left)
	step := new(big.Int).Div(gap, big.NewInt(int64(num)))
	remainder := new(big.Int).Mod(gap, big.NewInt(int64(num)))
	if step.Cmp(big.NewInt(0)) == 0 {
		num = int(remainder.Int64())
	}

	results = append(results, left)
	var lowerBound *big.Int
	upperBound := new(big.Int).Set(left)
	for i := 1; i < num; i++ {
		lowerBound = new(big.Int).Set(upperBound)
		upperBound = new(big.Int).Add(lowerBound, step)
		if remainder.Cmp(big.NewInt(int64(i))) >= 0 {
			upperBound = new(big.Int).Add(upperBound, big.NewInt(1))
		}
		results = append(results, upperBound)
	}
	results = append(results, right)
	return
}
