package runner

import (
	"os"

	mylog "github.com/Breeze0806/go/log"
)

var log mylog.Logger = mylog.NewDefaultLogger(os.Stderr, mylog.InfoLevel, "[datax]")

func init() {
	mylog.RegisterInitFuncs(func() {
		log = mylog.GetLogger()
	})
}
