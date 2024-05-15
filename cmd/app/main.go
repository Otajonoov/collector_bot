package main

import (
	"iman_tg_bot/internal/adapter"
	"iman_tg_bot/internal/pkg/config"
	db "iman_tg_bot/internal/pkg/db"
	"iman_tg_bot/internal/pkg/logger/slogpretty"
	"iman_tg_bot/internal/port/bot"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	//envDev   = "dev"
)

func main() {

	cfg := config.Load()

	log := setupLogger(cfg.Environment)

	pgxConn, err := db.ConnDB()
	if err != nil {
		log.Error("failed to connect to database", err)
		os.Exit(1)
	}
	defer pgxConn.Close()

	clientUser := adapter.NewRepo(pgxConn)

	bot := bot.NewBot(cfg, clientUser, log)

	bot.Start()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
