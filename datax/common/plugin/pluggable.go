package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

type Pluggable interface {
	Developer() (string, error)

	Description() (string, error)

	PluginName() (string, error)

	PluginJobConf() *config.JSON

	PeerPluginName() string

	PeerPluginJobConf() *config.JSON

	SetPluginJobConf(conf *config.JSON)

	SetPeerPluginJobConf(conf *config.JSON)

	SetPeerPluginName(name string)

	SetPluginConf(conf *config.JSON)

	Init(ctx context.Context) error

	Destroy(ctx context.Context) error
}

type BasePluggable struct {
	pluginConf        *config.JSON
	pluginJobConf     *config.JSON
	peerPluginName    string
	peerPluginJobConf *config.JSON
}

func NewBasePluggable() *BasePluggable {
	return &BasePluggable{}
}

func (b *BasePluggable) SetPluginConf(conf *config.JSON) {
	b.pluginConf = conf
}

func (b *BasePluggable) SetPluginJobConf(conf *config.JSON) {
	b.pluginJobConf = conf
}

func (b *BasePluggable) SetPeerPluginName(name string) {
	b.peerPluginName = name
}

func (b *BasePluggable) SetPeerPluginJobConf(conf *config.JSON) {
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

func (b *BasePluggable) PluginConf() *config.JSON {
	return b.pluginConf
}

func (b *BasePluggable) PluginJobConf() *config.JSON {
	return b.pluginJobConf
}

func (b *BasePluggable) PeerPluginName() string {
	return b.peerPluginName
}

func (b *BasePluggable) PeerPluginJobConf() *config.JSON {
	return b.peerPluginJobConf
}
