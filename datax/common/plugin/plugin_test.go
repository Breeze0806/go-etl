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
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func TestBasePlugin(t *testing.T) {
	ctx := context.Background()
	p := NewBasePlugin()
	if err := p.PreHandler(ctx, &config.JSON{}); err != nil {
		t.Errorf("PreHandler error  = %v", err)
	}
	if err := p.PreCheck(ctx); err != nil {
		t.Errorf("PreCheck error  = %v", err)
	}
	if err := p.Prepare(ctx); err != nil {
		t.Errorf("Prepare error  = %v", err)
	}
	if err := p.Post(ctx); err != nil {
		t.Errorf("Post error  = %v", err)
	}
	if err := p.PostHandler(ctx, &config.JSON{}); err != nil {
		t.Errorf("PostHandler error  = %v", err)
	}
}
