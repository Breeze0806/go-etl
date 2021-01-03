package job

import (
	"os"

	mylog "github.com/Breeze0806/go/log"
)

var log mylog.Logger = mylog.NewDefaultLogger(os.Stderr, mylog.ErrorLevel, "[datax]")

func init() {
	mylog.RegisterInitFuncs(func() {
		log = mylog.GetLogger()
	})
}
