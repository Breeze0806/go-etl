package taskgroup

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

type Option func(c *Container) error

type Container struct {
	jobID string
	conf  *config.JSON
	ctx   context.Context
}

func NewContainer(ctx context.Context, options ...Option) (c *Container, err error) {
	c = &Container{
		ctx: ctx,
	}
	for _, option := range options {
		if err = option(c); err != nil {
			return
		}
	}
	return
}

func setConfig(conf *config.JSON) Option {
	return func(c *Container) error {
		c.conf = conf
		return nil
	}
}

func setJobID(jobID string) Option {
	return func(c *Container) error {
		c.jobID = jobID
		return nil
	}
}

func (c *Container) Do() (err error) {
	return
}
