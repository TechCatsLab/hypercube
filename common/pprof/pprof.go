package pprof

import (
	_"net/http/pprof"
	"log"
	"net/http"
)

func HttpPprof()  {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

