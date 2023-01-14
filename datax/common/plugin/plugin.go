// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

//Plugin 插件
type Plugin interface {
	Pluggable
	//预检查
	PreCheck(ctx context.Context) error
	//准备
	Prepare(ctx context.Context) error
	//后置通知
	Post(ctx context.Context) error
	//预备处理， todo当前未用到
	PreHandler(ctx context.Context, conf *config.JSON) error
	//后置通知处理， todo当前未用到
	PostHandler(ctx context.Context, conf *config.JSON) error
}

//BasePlugin 基础插件，用于辅助和简化插件的实现
type BasePlugin struct {
	*BasePluggable
}

//NewBasePlugin 创建基础插件
func NewBasePlugin() *BasePlugin {
	return &BasePlugin{
		BasePluggable: NewBasePluggable(),
	}
}

//PreCheck 预检查空方法
func (b *BasePlugin) PreCheck(ctx context.Context) error {
	return nil
}

//Post 后置通知空方法
func (b *BasePlugin) Post(ctx context.Context) error {
	return nil
}

//Prepare 预备空方法
func (b *BasePlugin) Prepare(ctx context.Context) error {
	return nil
}

//PreHandler 预处理空方法
func (b *BasePlugin) PreHandler(ctx context.Context, conf *config.JSON) error {
	return nil
}

//PostHandler 后置通知处理空方法
func (b *BasePlugin) PostHandler(ctx context.Context, conf *config.JSON) error {
	return nil
}
