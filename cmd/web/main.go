package main

import (
	"flag"
	"os"

	"github.com/go-playground/form/v4"
	"github.com/rs/zerolog"
)

const version = "1.0.0"

type config struct {
	port    int
	env     string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}
type application struct {
	config      config
	logger      zerolog.Logger
	formDecoder *form.Decoder
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Application server port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development|production)")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	formDecoder := form.NewDecoder()
	app := &application{
		config:      cfg,
		logger:      logger,
		formDecoder: formDecoder,
	}

	logger.Info().Msg("Server starting")
	err := app.serve()

	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot start server")
	}

}
