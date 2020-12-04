package job

import "github.com/Breeze0806/go-etl/datax/core"

type Container struct {
	*core.BaseCotainer
}

func (c *Container) Start() error {
	return nil
}
