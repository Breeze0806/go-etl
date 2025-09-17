package mongodb

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
)

type Config struct {
	Uri        string `json:"uri"`
	Collection string `json:"collection"`
}

type Handler interface {
	GetConnection(ctx context.Context, conf *config.JSON) (*Client, error)
	//GetCollection(ctx context.Context, conf *config.JSON) (Collection, error)
}
type handler struct {
	conf *Config
}

func NewHandler() Handler {
	return &handler{}
}
func (h *handler) GetConnection(ctx context.Context, conf *config.JSON) (*Client, error) {
	uri, err := conf.GetString("uri")
	if err != nil {
		return nil, err
	}
	// Creates a new client and connects to the server
	client, err := NewClient(ctx, uri)
	if err != nil {
		return nil, err
	}
	// Sends a ping to confirm a successful connection
	err = client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}
