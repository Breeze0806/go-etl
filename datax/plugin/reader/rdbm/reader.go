package rdbm

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
)

type Reader interface {
	spi.Reader

	ResourcesConfig() *config.JSON
}

func RegisterReader(new func(string) (Reader, error)) (string, error) {
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
	loader.RegisterReader(name, writer)
	return pluginConfig, nil
}
