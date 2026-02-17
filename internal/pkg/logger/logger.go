package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(){

	// pastikan folder logs ada sebelum writer mencoba menulis file
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err){
		_ = os.Mkdir(logDir, 0755)
	}

	// buat writernya
	writer := &DailyWriter{
		LogDir: "logs",
	}

	// Masukkan writer kita ke handler slog
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
