package reader

import "github.com/Breeze0806/go-etl/datax/common/config"

type Job interface {
	Split(int) ([]*config.Json, error)
	
}
