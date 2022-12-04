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

package schedule

import (
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
	"github.com/pingcap/errors"
)

//NewRetryStrategy 根据配置文件生成重试策略
func NewRetryStrategy(j RetryJudger, conf *config.JSON) (s RetryStrategy, err error) {
	var retry *config.JSON
	if ok := conf.Exists("retry"); !ok {
		return NewNoneRetryStrategy(), nil
	}
	if retry, err = conf.GetConfig("retry"); err != nil {
		return
	}
	var typ string
	if typ, err = retry.GetString("type"); err != nil {
		return
	}

	var strategy *config.JSON
	if strategy, err = retry.GetConfig("strategy"); err != nil {
		return
	}

	switch typ {
	case "ntimes":
		var retryConf NTimesRetryConfig
		if err = json.Unmarshal([]byte(strategy.String()), &retryConf); err != nil {
			return
		}
		if retryConf.N == 0 || retryConf.Wait.Duration == 0 {
			err = errors.New("ntimes retry config is valid")
			return
		}
		s = NewNTimesRetryStrategy(j, retryConf.N, retryConf.Wait.Duration)
		return
	case "forever":
		var retryConf ForeverRetryConfig
		if err = json.Unmarshal([]byte(strategy.String()), &retryConf); err != nil {
			return
		}
		if retryConf.Wait.Duration == 0 {
			err = errors.New("forever retry config is valid")
			return
		}
		s = NewForeverRetryStrategy(j, retryConf.Wait.Duration)
		return
	case "exponential":
		var retryConf ExponentialRetryConfig
		if err = json.Unmarshal([]byte(strategy.String()), &retryConf); err != nil {
			return
		}
		if retryConf.Init.Duration == 0 || retryConf.Max.Duration == 0 {
			err = errors.New("exponential retry config is valid")
			return
		}
		s = NewExponentialRetryStrategy(j, retryConf.Init.Duration, retryConf.Max.Duration)
		return
	}
	err = errors.Errorf("no such type(%v)", typ)
	return
}

//NTimesRetryConfig n次数重复重试策略
type NTimesRetryConfig struct {
	N    int            `json:"n"`
	Wait time2.Duration `json:"wait"`
}

//ForeverRetryConfig 永久重复重试策略
type ForeverRetryConfig struct {
	Wait time2.Duration `json:"wait"`
}

//ExponentialRetryConfig 幂重复重试策略
type ExponentialRetryConfig struct {
	Init time2.Duration `json:"init"`
	Max  time2.Duration `json:"max"`
}

//RetryStrategy 重试策略
type RetryStrategy interface {
	Next(err error, n int) (retry bool, wait time.Duration)
}

//RetryJudger 重试判断器
type RetryJudger interface {
	ShouldRetry(err error) bool
}

//NoneRetryStrategy 无重试策略
type NoneRetryStrategy struct{}

//NewNoneRetryStrategy 创建无重试策略
func NewNoneRetryStrategy() RetryStrategy {
	return &NoneRetryStrategy{}
}

//Next 下一次是否retry需要重试，wait等待时间
func (r *NoneRetryStrategy) Next(err error, n int) (retry bool, wait time.Duration) {
	return
}

//NTimesRetryStrategy n次数重复重试策略
type NTimesRetryStrategy struct {
	j    RetryJudger
	n    int
	wait time.Duration
}

//NewNTimesRetryStrategy 通过重试判定器j,最大次数n以及重试间隔wait创建n次数重复重试策略
func NewNTimesRetryStrategy(j RetryJudger, n int, wait time.Duration) RetryStrategy {
	return &NTimesRetryStrategy{
		j:    j,
		n:    n,
		wait: wait,
	}
}

//Next 通过错误err以及当前次数n获取下次是否重试retry以及下次时间间隔wait
func (r *NTimesRetryStrategy) Next(err error, n int) (retry bool, wait time.Duration) {
	if !r.j.ShouldRetry(err) {
		return false, 0
	}

	if n >= r.n {
		return false, 0
	}
	return true, r.wait
}

//ForeverRetryStrategy 永久重试策略
type ForeverRetryStrategy struct {
	j    RetryJudger
	wait time.Duration
}

//NewForeverRetryStrategy 通过重试判定器j以及重试间隔wait创建永久重试策略
func NewForeverRetryStrategy(j RetryJudger, wait time.Duration) RetryStrategy {
	return &ForeverRetryStrategy{
		j:    j,
		wait: wait,
	}
}

//Next 通过错误err,获取下次是否重试retry以及下次时间间隔wait,在永久重试策略没有最大重试次数，当前次数n没有作用
func (r *ForeverRetryStrategy) Next(err error, _ int) (retry bool, wait time.Duration) {
	if !r.j.ShouldRetry(err) {
		return false, 0
	}

	return true, r.wait
}

//ExponentialStrategy 幂重试策略
type ExponentialStrategy struct {
	j    RetryJudger
	f    float64
	init float64
	max  float64
}

//NewExponentialRetryStrategy 通过重试判定器j,开始时间间隔init以及最大时间间隔max创建幂重试策略
func NewExponentialRetryStrategy(j RetryJudger, init, max time.Duration) RetryStrategy {
	rand.Seed(time.Now().UnixNano())
	return &ExponentialStrategy{
		j:    j,
		f:    2.0,
		init: float64(init),
		max:  float64(max),
	}
}

//Next 通过错误err,获取下次是否重试retry以及下次时间间隔wait,在幂重试策略最大时间间隔
func (r *ExponentialStrategy) Next(err error, n int) (retry bool, wait time.Duration) {
	if !r.j.ShouldRetry(err) {
		return false, 0
	}
	x := 1.0 + rand.Float64() // random number in [1..2]
	m := math.Min(x*r.init*math.Pow(r.f, float64(n)), r.max)
	if m >= r.max {
		return false, 0
	}
	return true, time.Duration(m)
}
