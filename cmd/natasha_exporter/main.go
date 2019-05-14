package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kaminek/natasha_exporter/pkg/config"
	"github.com/kaminek/natasha_exporter/pkg/exporter"
	"github.com/kaminek/natasha_exporter/pkg/info"
	"github.com/kaminek/natasha_exporter/pkg/server"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	cfg := config.Load()

	if env := os.Getenv("NATASHA_EXPORTER_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:    "Natasha_exporter",
		Version: info.Version,
		Usage:   "Natasha Exporter",
		Authors: []*cli.Author{
			{
				Name:  "Amine KHERBOUCHE",
				Email: "akherbouche@scaleway.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "web.address",
				Value:       "0.0.0.0:9507",
				Usage:       "Address to bind the metrics server",
				EnvVars:     []string{"NATASHA_EXPORTER_WEB_ADDRESS"},
				Destination: &cfg.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "server.entry",
				Value:       "/natasha_metrics",
				Usage:       "Metrics path",
				EnvVars:     []string{"NATASHA_EXPORTER_PATH"},
				Destination: &cfg.Server.Path,
			},
			&cli.DurationFlag{
				Name:        "natasha.timeout",
				Value:       5 * time.Second,
				Usage:       "Target request timeout as duration",
				EnvVars:     []string{"NATASHA_EXPORTER_TIMEOUT"},
				Destination: &cfg.Target.Timeout,
			},
		},
		Action: func(c *cli.Context) error {
			exporter.New(cfg)
			return server.Init(cfg)
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
