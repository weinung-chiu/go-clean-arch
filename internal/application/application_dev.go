package application

import (
	"context"
	"log/slog"
)

// DevelopApplication 用來開發
type DevelopApplication struct {
	logger *slog.Logger
}

type NewDevAppParams struct {
	Logger *slog.Logger
}

func NewDevelopApplication(params NewDevAppParams) *DevelopApplication {
	return &DevelopApplication{
		// 初始化時加上 component ，寫入 record 時即可表達是由誰寫入的
		logger: params.Logger.With("component", "DevelopApplication"),
	}
}

func (app *DevelopApplication) DoSomething(ctx context.Context) error {
	app.logger.InfoContext(ctx, "did something")

	return nil
}

func (app *DevelopApplication) DoSomethingFatal(ctx context.Context) error {
	app.logger.WarnContext(ctx, "did something fatal", "foo", "bar")

	return nil
}
