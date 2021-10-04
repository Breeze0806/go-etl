package main

import (
	"os"

	mylog "github.com/Breeze0806/go/log"
)

var log = mylog.NewDefaultLogger(os.Stdout, mylog.DebugLevel, "[datax]")

func init() {
	f, err := os.Create("datax.log")
	if err != nil {
		panic(err)
	}
	log = mylog.NewDefaultLogger(f, mylog.DebugLevel, "[datax]")
}

func initLog() {
	mylog.SetLogger(log)
}
