package application

import (
	"context"
	"log/slog"
)

// ExampleApplication 用來開發
type ExampleApplication struct {
	logger *slog.Logger
}

type NewExampleAppParams struct {
	Logger *slog.Logger
}

func NewExampleApplication(params NewExampleAppParams) *ExampleApplication {
	return &ExampleApplication{
		// 初始化時加上 component ，寫入 record 時即可表達是由誰寫入的
		logger: params.Logger.With("component", "ExampleApplication"),
	}
}

func (app *ExampleApplication) DoSomething(ctx context.Context) error {
	app.logger.InfoContext(ctx, "did something")

	return nil
}

func (app *ExampleApplication) DoSomethingFatal(ctx context.Context) error {
	app.logger.WarnContext(ctx, "did something fatal", "foo", "bar")

	return nil
}
