package runner

import mylog "github.com/Breeze0806/go/log"

var log mylog.Logger

func init() {
	mylog.RegisterInitFuncs(func() {
		log = mylog.GetLogger()
	})
}
