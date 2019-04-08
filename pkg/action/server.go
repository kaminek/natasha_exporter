package action

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kaminek/natasha_exporter/pkg/config"
	"github.com/kaminek/natasha_exporter/pkg/exporter"
	"github.com/kaminek/natasha_exporter/pkg/version"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server handles the server sub-command.
func Server(cfg *config.Config, logger log.Logger) error {
	level.Info(logger).Log(
		"msg", "Launching Natasha Exporter",
		"version", version.Version,
		"revision", version.Revision,
		"date", version.BuildDate,
		"go", version.GoVersion,
	)

	var gr run.Group

	{
		server := &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      handler(cfg, logger),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "Starting metrics server",
				"addr", cfg.Server.Addr,
			)

			return server.ListenAndServe()
		}, func(reason error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				level.Error(logger).Log(
					"msg", "Failed to shutdown metrics gracefully",
					"err", err,
				)

				return
			}

			level.Info(logger).Log(
				"msg", "Metrics shutdown gracefully",
				"reason", reason,
			)
		})
	}

	{
		stop := make(chan os.Signal, 1)

		gr.Add(func() error {
			signal.Notify(stop, os.Interrupt)

			<-stop

			return nil
		}, func(err error) {
			close(stop)
		})
	}

	return gr.Run()
}

func handler(cfg *config.Config, logger log.Logger) *chi.Mux {
	mux := chi.NewRouter()

	r := prometheus.NewRegistry()
	// r.MustRegister(prometheus.NewProcessCollector(nil))
	r.MustRegister(prometheus.NewGoCollector())

	r.MustRegister(exporter.NewGeneralCollector(
		version.Version,
		version.Revision,
		version.BuildDate,
		version.GoVersion,
		version.StartTime,
	))

	requestFailures := exporter.RequestFailures()
	r.MustRegister(requestFailures)

	requestDuration := exporter.RequestDuration()
	r.MustRegister(requestDuration)

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.Server.Path, http.StatusMovedPermanently)
	})

	mux.Route("/", func(root chi.Router) {
		root.Mount(
			cfg.Server.Path,
			promhttp.HandlerFor(
				r,
				promhttp.HandlerOpts{},
			),
		)

		root.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})
	})

	return mux
}
