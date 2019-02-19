package main

import (
	"flag"
	"log"

	"github.com/jserver/floader"
	"github.com/jserver/server"
)

func main() {
	var (
		addr  = flag.String("http", ":8080", "HTTP service address (e.g., ':8080')")
		jfile = flag.String("file", "jserver.json", "json file that to load")
	)
	flag.Parse()
	log.Printf("Listen [%s],Using [%s]\n", *addr, *jfile)
	fl := floader.NewFloader(*jfile)
	go fl.Run()

	var srv *server.Server
	for {
		data := fl.Sync()
		srv.Stop()
		srv = server.NewServer(*addr, data)
		srv.Start()
	}

}
