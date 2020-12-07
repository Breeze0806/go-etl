package plugin

import "github.com/Breeze0806/go-etl/datax/common/config"

type Plugin interface {
	Pluggable
	PreCheck() error
	Prepare() error
	Post() error
	PreHandler(conf *config.Json) error
	PostHandler(conf *config.Json) error
}

type BasePlugin struct {
	*BasePluggable
}

func (b *BasePlugin) PreCheck() error {
	return nil
}

func (b *BasePlugin) Post() error {
	return nil
}

func (b *BasePlugin) Prepare() error {
	return nil
}

func (b *BasePlugin) PreHandler(conf *config.Json) error {
	return nil
}

func (b *BasePlugin) PostHandler(conf *config.Json) error {
	return nil
}
