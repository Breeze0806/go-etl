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
	"fmt"
	"os"

	"github.com/Breeze0806/go-etl/cmd/datax/tools"
)

func main() {
	initLog()
	configFile := flag.String("c", "F:\\OpenSource\\etl\\go-etl\\cmd\\datax\\config.json", "config")
	wizardFile := flag.String("w", "", "wizard")
	httpAddr := flag.String("http", "", "http")
	flag.Parse()
	if *wizardFile != "" {
		if err := tools.NewWizard(*configFile, *wizardFile).GenerateConfigsAndScripts(); err != nil {
			fmt.Printf("wizard generate configs fail. err: %v\n", err)
			log.Errorf("wizard generate configs fail. err: %v", err)
			os.Exit(1)
		}
		return
	}

	log.Infof("config: %v\n", *configFile)

	e := newEnveronment(*configFile, *httpAddr)
	defer e.close()
	if err := e.build(); err != nil {
		fmt.Printf("run fail. err: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("run success\n")
}
