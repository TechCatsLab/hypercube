package main

import (
	_"net/http/pprof"
	"log"
	"net/http"
)

func HttpPprof()  {
	go func() {
		log.Println(http.ListenAndServe(configuration.PprofAddrs, nil))
	}()
}
