package writer

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
)

//Writer 写入器
type Writer interface {
	spi.Writer

	//资源插件配置
	ResourcesConfig() *config.JSON
}

type WriterMaker interface {
	Default() (Writer, error)
	FromFile(filename string) (Writer, error)
}

//RegisterWriter 注册创建函数new写入器,返回的是资源插件配置文件地州，出错时会返回error
func RegisterWriter(maker WriterMaker) (pluginConfig string, err error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("fail to get filename")
	}
	path := filepath.Dir(file)
	pluginConfig = filepath.Join(path, "resources", "plugin.json")
	var writer Writer
	if writer, err = maker.FromFile(pluginConfig); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		if writer, err = maker.Default(); err != nil {
			return "", err
		}
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
