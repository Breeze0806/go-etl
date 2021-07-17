package rdbm

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
)

type Writer interface {
	spi.Writer

	ResourcesConfig() *config.JSON
}

func RegisterWriter(new func(string) (Writer, error)) (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("fail to get filename")
	}
	path := filepath.Dir(file)
	pluginConfig := filepath.Join(path, "resources", "plugin.json")
	writer, err := new(pluginConfig)
	if err != nil {
		return "", err
	}
	name, err := writer.ResourcesConfig().GetString("name")
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", errors.New("name is empty")
	}
	loader.RegisterWriter(name, writer)
	return pluginConfig, nil
}
