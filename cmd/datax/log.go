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

package main

import (
	"os"

	mylog "github.com/Breeze0806/go/log"
)

var log = mylog.NewDefaultLogger(os.Stdout, mylog.DebugLevel, "[go-etl]")

func init() {
	f, err := os.OpenFile("go-etl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log = mylog.NewDefaultLogger(f, mylog.DebugLevel, "[go-etl]")
}

func initLog() {
	mylog.SetLogger(log)
}
