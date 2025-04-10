package main

import (
	"context"
	"log/slog"
	"os"

	"go-clean-arch/internal/application"
	"go-clean-arch/internal/common"
)

var traceIDKey = "traceIDKey"

func main() {
	rootLogger := slog.New(common.NewTraceHandler(slog.NewJSONHandler(os.Stdout, nil)))
	//rootLogger := slog.New(common.NewTraceHandler(slog.NewTextHandler(os.Stdout, nil)))
	rootLogger.Info("initial app with log/slog", "my_key", "my_value")

	rootCtx := context.Background()

	// logger 在 constructor 顯式帶入
	params := application.NewDevAppParams{Logger: rootLogger}
	app := application.NewDevelopApplication(params)

	_ = app.DoSomething(rootCtx)
	_ = app.DoSomethingFatal(rootCtx)
}
