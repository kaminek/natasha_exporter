package server

import (
	"net/http"

	"github.com/kaminek/natasha_exporter/pkg/config"
	"github.com/kaminek/natasha_exporter/pkg/exporter"
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
             <head><title>Haproxy Exporter</title></head>
             <body>
             <h1>Haproxy Exporter</h1>
             <p><a href='` + cfg.Server.Path + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, nil))
}
