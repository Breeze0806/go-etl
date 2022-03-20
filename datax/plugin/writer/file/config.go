package file

import (
	"encoding/json"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

var (
	defalutBatchSize    = 1000
	defalutBatchTimeout = 1 * time.Second
)

type Config interface {
	GetBatchSize() int              //单次批量写入数
	GetBatchTimeout() time.Duration //单次批量写入超时时间
}

type BaseConfig struct {
	BatchSize    int            `json:"batchSize"`
	BatchTimeout time2.Duration `json:"batchTimeout"`
}

func NewBaseConfig(conf *config.JSON) (*BaseConfig, error) {
	c := &BaseConfig{}
	if err := json.Unmarshal([]byte(conf.String()), c); err != nil {
		return nil, err
	}
	return c, nil
}

func (b *BaseConfig) GetBatchTimeout() time.Duration {
	if b.BatchTimeout.Duration == 0 {
		return defalutBatchTimeout
	}
	return b.BatchTimeout.Duration
}

func (b *BaseConfig) GetBatchSize() int {
	if b.BatchSize == 0 {
		return defalutBatchSize
	}

	return b.BatchSize
}
