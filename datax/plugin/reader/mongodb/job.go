package mongodb

import (
	"context"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/pkg/errors"
)

// Job normal dbms job
type Job struct {
	*plugin.BaseJob
	handler Handler
	client  *Client
	uri     string
}

func NewJob(handler Handler) *Job {
	return &Job{
		BaseJob: plugin.NewBaseJob(),
		handler: handler,
	}
}

func (j *Job) Init(ctx context.Context) error {
	// test connection
	connection, err := j.PluginJobConf().GetConfig("connection")
	if err != nil {
		return errors.Wrap(err, "GetConfig connection fail")
	}
	c, err := config.NewJSONFromString("{}")
	if err != nil {
		return errors.Wrap(err, "NewJSONFromString fail")
	}
	address, err := connection.GetString("address")
	if err != nil {
		return errors.Wrap(err, "GetString address fail")
	}
	username, err := connection.GetString("username")
	if err != nil {
		return errors.Wrap(err, "GetString username fail")
	}
	password, err := connection.GetString("password")
	if err != nil {
		return errors.Wrap(err, "GetString password fail")
	}
	uri := fmt.Sprintf("mongodb://%s", address)
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", username, password, address)
	}
	c.Set("uri", uri)
	j.uri = uri
	j.client, err = j.handler.GetConnection(ctx, c)
	if err != nil {
		return errors.Wrap(err, "TestConnection fail")
	}
	return nil
}
func (j *Job) Destroy(ctx context.Context) error {
	return nil
}

// Split - 根据作业配置将任务分割为多个子任务
func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	// test connection
	connection, err := j.PluginJobConf().GetConfig("connection")
	if err != nil {
		return nil, errors.Wrap(err, "GetConfig connection fail")
	}
	table, err := connection.GetConfig("table")
	if err != nil {
		return nil, errors.Wrap(err, "GetConfig table fail")
	}
	database, err := table.GetString("db")
	if err != nil {
		return nil, errors.Wrap(err, "GetString database fail")
	}
	collection, err := table.GetString("collection")
	if err != nil {
		return nil, errors.Wrap(err, "GetString collection fail")
	}
	idRange, err := j.client.GetObjectIDRange(ctx, database, collection, "_id", int64(number))
	if err != nil {
		return nil, errors.Wrap(err, "GetIdRange fail")
	}

	configs := make([]*config.JSON, number)

	taskConf, _ := config.NewJSONFromString("{}")
	taskConf.Set("database", database)
	taskConf.Set("collection", collection)
	taskConf.Set("uri", j.uri)
	for i := 0; i < len(idRange); i++ {
		tkConf := taskConf.CloneConfig()
		tkConf.Set("max", idRange[i].Max)
		tkConf.Set("min", idRange[i].Min)
		tkConf.Set("taskId", fmt.Sprintf("task_%d", i))
		configs[i] = tkConf
	}
	fmt.Println(configs)
	return configs, nil
}
