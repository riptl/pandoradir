package main

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	original   http.ResponseWriter
	path       string
	remoteAddr string
	executed   bool
}

func (l loggingResponseWriter) Header() http.Header {
	return l.original.Header()
}

func (l loggingResponseWriter) Write(x []byte) (int, error) {
	return l.original.Write(x)
}

func (l loggingResponseWriter) WriteHeader(sc int) {
	l.original.Header().Set("Server", "pandoradir")
	log.Printf("%s: %d %s", l.remoteAddr, sc, l.path)
	l.executed = true
	l.original.WriteHeader(sc)
}

func (l loggingResponseWriter) Finalize() {
	if !l.executed {
		log.Printf("%s: 200 %s", l.remoteAddr, l.path)
	}
}

func logHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		lw := loggingResponseWriter{
			original:   res,
			path:       req.URL.Path,
			remoteAddr: req.RemoteAddr,
		}
		f(lw, req)
		lw.Finalize()
	}
}
