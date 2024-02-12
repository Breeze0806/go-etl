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

// Plugin: an extension or add-on component
type Plugin interface {
	Pluggable
	// PreCheck: a pre-processing check or verification step
	PreCheck(ctx context.Context) error
	// Prepare: a preparation step before the main operation
	Prepare(ctx context.Context) error
	// PostNotification: a notification step after the main operation
	Post(ctx context.Context) error
	// PreHandler: preprocessing, todo: currently not in use
	PreHandler(ctx context.Context, conf *config.JSON) error
	// PostHandler: post-notification processing, todo: currently not in use
	PostHandler(ctx context.Context, conf *config.JSON) error
}

// BasePlugin: a fundamental plugin class that assists and simplifies the implementation of plugins
type BasePlugin struct {
	*BasePluggable
}

// NewBasePlugin: a function or method to create a new instance of BasePlugin
func NewBasePlugin() *BasePlugin {
	return &BasePlugin{
		BasePluggable: NewBasePluggable(),
	}
}

// PreCheck: an empty method for pre-checking
func (b *BasePlugin) PreCheck(ctx context.Context) error {
	return nil
}

// Post: an empty method for post-notification
func (b *BasePlugin) Post(ctx context.Context) error {
	return nil
}

// Prepare: an empty method for preparation
func (b *BasePlugin) Prepare(ctx context.Context) error {
	return nil
}

// PreHandler: an empty method for preprocessing
func (b *BasePlugin) PreHandler(ctx context.Context, conf *config.JSON) error {
	return nil
}

// PostHandler: an empty method for post-notification processing
func (b *BasePlugin) PostHandler(ctx context.Context, conf *config.JSON) error {
	return nil
}
