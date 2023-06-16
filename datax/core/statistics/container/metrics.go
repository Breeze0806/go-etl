package container

import (
	"sync"

	"github.com/Breeze0806/go/encoding"
)

//Metrics json格式指标
type Metrics struct {
	sync.RWMutex

	metricJSON *encoding.JSON
}

//NewMetrics json格式指标
func NewMetrics() *Metrics {
	j, _ := encoding.NewJSONFromString("{}")
	return &Metrics{
		metricJSON: j,
	}
}

//JSON json格式指标
func (m *Metrics) JSON() *encoding.JSON {
	m.RLock()
	defer m.RUnlock()
	return m.metricJSON
}

//Set 设置path的value
func (m *Metrics) Set(path string, value interface{}) error {
	m.Lock()
	defer m.Unlock()
	return m.metricJSON.Set(path, value)
}

//Get 获得path的value
func (m *Metrics) Get(path string) *encoding.JSON {
	m.RLock()
	defer m.RUnlock()
	j, err := m.metricJSON.GetJSON(path)
	if err != nil {
		return nil
	}
	return j
}
