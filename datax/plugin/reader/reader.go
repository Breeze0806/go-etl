package reader

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi"
)

//Reader 数据库读取器
type Reader interface {
	spi.Reader

	ResourcesConfig() *config.JSON //插件资源配置
}

type ReaderMaker interface {
	Default() (Reader, error)
	FromFile(filename string) (Reader, error)
}

//RegisterReader 通过生成数据库读取器函数new注册到读取器，返回插件资源配置文件地址，在出错时会包err
//目前未在代码中实际使用，而是通过datax/build的go generate命令自动将resources/plugin.json
//中的内容放入到新生成的代码文件中，用以注册Reader
func RegisterReader(maker ReaderMaker) (pluginConfig string, err error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("fail to get filename")
	}
	path := filepath.Dir(file)
	pluginConfig = filepath.Join(path, "resources", "plugin.json")
	var reader Reader

	if reader, err = maker.FromFile(pluginConfig); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		if reader, err = maker.Default(); err != nil {
			return "", err
		}
	}
	name := ""
	name, err = reader.ResourcesConfig().GetString("name")
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", errors.New("name is empty")
	}
	loader.RegisterReader(name, reader)
	return pluginConfig, nil
}
