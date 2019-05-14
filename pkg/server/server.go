package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/kaminek/natasha_exporter/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

// Init HTTP Server.
func Init(cfg *config.Config) error {
	log.Infoln("Starting Natasha Exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle(cfg.Server.Path, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Natasha Exporter</title></head>
             <body>
             <h1>Natasha Exporter</h1>
             <p><a href='` + cfg.Server.Path + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, nil))
	return nil
}

// NatashaServerDial talks to server and fills struct
func NatashaServerDial() (net.Conn, error) {

	conn, err := net.Dial("tcp4", "localhost:4242")
	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem while connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		return nil, err
	}
	return conn, err
}
