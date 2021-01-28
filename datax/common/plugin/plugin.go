package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

//Plugin 插件
type Plugin interface {
	Pluggable
	PreCheck(ctx context.Context) error
	Prepare(ctx context.Context) error
	Post(ctx context.Context) error
	PreHandler(ctx context.Context, conf *config.JSON) error
	PostHandler(ctx context.Context, conf *config.JSON) error
}

type BasePlugin struct {
	*BasePluggable
}

func NewBasePlugin() *BasePlugin {
	return &BasePlugin{
		BasePluggable: NewBasePluggable(),
	}
}

func (b *BasePlugin) PreCheck(ctx context.Context) error {
	return nil
}

func (b *BasePlugin) Post(ctx context.Context) error {
	return nil
}

func (b *BasePlugin) Prepare(ctx context.Context) error {
	return nil
}

func (b *BasePlugin) PreHandler(ctx context.Context, conf *config.JSON) error {
	return nil
}

func (b *BasePlugin) PostHandler(ctx context.Context, conf *config.JSON) error {
	return nil
}
