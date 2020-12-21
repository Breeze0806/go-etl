package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/config"
)

type Pluggable interface {
	Developer() (string, error)

	Description() (string, error)

	PluginName() (string, error)

	PluginJobConf() *config.Json

	PeerPluginName() string

	PeerPluginJobConf() *config.Json

	SetPluginJobConf(conf *config.Json)

	SetPeerPluginJobConf(conf *config.Json)

	SetPeerPluginName(name string)

	SetPluginConf(conf *config.Json)

	Init(ctx context.Context) error

	Destroy(ctx context.Context) error
}

type BasePluggable struct {
	pluginConf        *config.Json
	pluginJobConf     *config.Json
	peerPluginName    string
	peerPluginJobConf *config.Json
}

func NewBasePluggable() *BasePluggable {
	return &BasePluggable{}
}

func (b *BasePluggable) SetPluginConf(conf *config.Json) {
	b.pluginConf = conf
}

func (b *BasePluggable) SetPluginJobConf(conf *config.Json) {
	b.pluginJobConf = conf
}

func (b *BasePluggable) SetPeerPluginName(name string) {
	b.peerPluginName = name
}

func (b *BasePluggable) SetPeerPluginJobConf(conf *config.Json) {
	b.peerPluginJobConf = conf
}

func (b *BasePluggable) Developer() (string, error) {
	return b.pluginConf.GetString("developer")
}

func (b *BasePluggable) Description() (string, error) {
	return b.pluginConf.GetString("description")
}

func (b *BasePluggable) PluginName() (string, error) {
	return b.pluginConf.GetString("name")
}

func (b *BasePluggable) PluginConf() *config.Json {
	return b.pluginConf
}

func (b *BasePluggable) PluginJobConf() *config.Json {
	return b.pluginJobConf
}

func (b *BasePluggable) PeerPluginName() string {
	return b.peerPluginName
}

func (b *BasePluggable) PeerPluginJobConf() *config.Json {
	return b.peerPluginJobConf
}
