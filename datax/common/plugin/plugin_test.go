package plugin

import (
	"context"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/config"
)

func TestBasePlugin(t *testing.T) {
	ctx := context.Background()
	p := NewBasePlugin()
	if err := p.PreHandler(ctx, &config.Json{}); err != nil {
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
	if err := p.PostHandler(ctx, &config.Json{}); err != nil {
		t.Errorf("PostHandler error  = %v", err)
	}
}
