package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Column 表示 MongoDB 字段配置
type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Spliter string `json:"spliter,omitempty"`
}

type Task struct {
	*plugin.BaseTask
	client  *Client
	columns []Column
}

func NewTask() *Task {
	return &Task{
		BaseTask: plugin.NewBaseTask(),
	}
}

func (t *Task) Init(ctx context.Context) (err error) {
	// test connection
	conf := t.PluginJobConf()
	uri, err := conf.GetString("uri")
	if err != nil {
		return errors.Wrap(err, "uri not found")
	}
	t.client, err = NewClient(ctx, uri)
	if err != nil {
		return errors.Wrap(err, "create client failed")
	}

	// 解析 column 配置
	columnConf, err := conf.GetConfig("column")
	if err == nil {
		var columns []Column
		if err := json.Unmarshal([]byte(columnConf.String()), &columns); err != nil {
			return errors.Wrap(err, "parse column config failed")
		}
		t.columns = columns
	}

	fmt.Println(conf)

	return nil
}

func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.client != nil {
		t.client.Close()
	}
	return nil
}

func (t *Task) StartRead(ctx context.Context, p plugin.RecordSender) (err error) {
	fmt.Println("start read")
	conf := t.PluginJobConf()
	maxID, err := conf.GetString("max")

	if err != nil {
		return errors.Wrap(err, "max not found")
	}
	fmt.Println(maxID)
	minID, err := conf.GetString("min")
	if err != nil {
		return errors.Wrap(err, "min not found")
	}
	database, err := conf.GetString("database")
	if err != nil {
		return errors.Wrap(err, "database not found")
	}
	collection, err := conf.GetString("collection")
	if err != nil {
		return errors.Wrap(err, "collection not found")
	}
	key, err := conf.GetString("split_key")
	if err != nil {
		key = "_id"
	}

	fmt.Println(database, collection)
	data, err := t.client.GetDocByRange(ctx, database, collection, key, minID, maxID)
	// 遍历查询结果并转换为 Record
	for _, doc := range data {
		fmt.Println(doc)
		// 创建一个新的 Record
		record, err := p.CreateRecord()
		if err != nil {
			return errors.Wrap(err, "create record failed")
		}

		// 将 BSON 文档转换为 Record 的列
		err = t.convertBSONToRecord(doc, record)
		if err != nil {
			return errors.Wrap(err, "convert bson to record failed")
		}

		// 发送 Record
		err = p.SendWriter(record)
		if err != nil {
			return errors.Wrap(err, "send record failed")
		}
	}

	// 终止发送
	p.Terminate()
	return nil
}

// convertBSONToRecord 将 BSON 文档转换为 DataX Record
func (t *Task) convertBSONToRecord(doc bson.M, record element.Record) error {
	// 如果配置了 column，则只处理配置中的字段
	if len(t.columns) > 0 {
		for _, col := range t.columns {
			value, exists := doc[col.Name]
			if !exists {
				// 如果字段不存在，添加一个空值
				column := element.NewDefaultColumn(element.NewNilStringColumnValue(), col.Name, 0)
				record.Add(column)
				continue
			}

			column, err := t.convertValueToColumn(col.Name, value)
			if err != nil {
				return errors.Wrapf(err, "convert field %s failed", col.Name)
			}
			record.Add(column)
		}
	} else {
		// 如果没有配置 column，处理所有字段
		for key, value := range doc {
			column, err := t.convertValueToColumn(key, value)
			if err != nil {
				return errors.Wrapf(err, "convert field %s failed", key)
			}
			record.Add(column)
		}
	}
	return nil
}

// convertValueToColumn 将 BSON 值转换为 DataX Column
func (t *Task) convertValueToColumn(key string, value interface{}) (element.Column, error) {
	var columnValue element.ColumnValue

	// 查找字段配置
	var colConfig *Column
	for _, col := range t.columns {
		if col.Name == key {
			colConfig = &col
			break
		}
	}

	switch v := value.(type) {
	case string:
		columnValue = element.NewStringColumnValue(v)
	case int:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(fmt.Sprintf("%d", v))
		} else {
			columnValue = element.NewBigIntColumnValueFromInt64(int64(v))
		}
	case int32:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(fmt.Sprintf("%d", v))
		} else {
			columnValue = element.NewBigIntColumnValueFromInt64(int64(v))
		}
	case int64:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(fmt.Sprintf("%d", v))
		} else {
			columnValue = element.NewBigIntColumnValueFromInt64(v)
		}
	case float32:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(fmt.Sprintf("%g", v))
		} else {
			columnValue = element.NewDecimalColumnValueFromFloat(float64(v))
		}
	case float64:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(fmt.Sprintf("%g", v))
		} else {
			columnValue = element.NewDecimalColumnValueFromFloat(v)
		}
	case bool:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(fmt.Sprintf("%t", v))
		} else {
			columnValue = element.NewBoolColumnValue(v)
		}
	case nil:
		columnValue = element.NewNilStringColumnValue()
	case primitive.ObjectID:
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(v.Hex())
		} else {
			// ObjectID默认作为字符串处理
			columnValue = element.NewStringColumnValue(v.Hex())
		}
	case primitive.DateTime:
		tm := v.Time()
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(tm.Format("2006-01-02 15:04:05"))
		} else {
			columnValue = element.NewTimeColumnValue(tm)
		}
	case []interface{}:
		// 处理数组类型
		if colConfig != nil && colConfig.Type == "Array" {
			// 根据 spliter 配置处理数组
			strValues := make([]string, len(v))
			for i, item := range v {
				strValues[i] = fmt.Sprintf("%v", item)
			}
			spliter := " " // 默认分隔符为空格
			if colConfig.Spliter != "" {
				spliter = colConfig.Spliter
			}
			columnValue = element.NewStringColumnValue(strings.Join(strValues, spliter))
		} else {
			// 默认序列化为 JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			columnValue = element.NewStringColumnValue(string(jsonBytes))
		}
	case bson.M:
		// 处理嵌套文档
		if colConfig != nil && colConfig.Type == "string" {
			// 如果配置为字符串类型，则转换为 JSON 字符串
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			columnValue = element.NewStringColumnValue(string(jsonBytes))
		} else {
			// 默认序列化为 JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			columnValue = element.NewStringColumnValue(string(jsonBytes))
		}
	default:
		// 对于其他类型，尝试转换为字符串
		strValue := fmt.Sprintf("%v", v)
		if colConfig != nil && colConfig.Type == "string" {
			columnValue = element.NewStringColumnValue(strValue)
		} else {
			columnValue = element.NewStringColumnValue(strValue)
		}
	}

	// 使用 NewDefaultColumn 创建 Column，提供列名和估计的字节大小
	return element.NewDefaultColumn(columnValue, key, 0), nil
}
