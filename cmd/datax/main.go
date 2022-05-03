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
	"flag"
	"os"
)

func main() {
	initLog()
	var filename = flag.String("c", "config.json", "config")
	flag.Parse()
	log.Debugf("%v", *filename)
	e := newEnveronment(*filename)
	defer e.close()
	if err := e.build(); err != nil {
		log.Errorf("build fail. err : %v", err)
		os.Exit(1)
	}
	return
}
