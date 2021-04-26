package main

import (
	"flag"
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
		return
	}
	return
}
