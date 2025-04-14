package application

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"go-clean-arch/internal/common/random"
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

type RandomExample struct {
	MockRNG   int
	PseudoRNG int
	CryptoRNG int
}

// DemoRandom 示範三種 RNG 的結果
//
// Mock 為手動指定：一定會選中權重 2
// pseudo 指定固定的 source (種子) 324，偽隨機結果為 10076，一定會選中權重 1
// crypto 為真隨機，最有可能會選中權重 0
//
// 合理來說應該在 main.go injection 到 application 內而不是在 function 內初始化
// 但視需求也可以如此 Example 般全部使用。
func (app *ExampleApplication) DemoRandom(ctx context.Context) (RandomExample, error) {
	weights := []int{10000, 100, 1}

	mockRNG := &random.MockRNG{NextInt: 10101}
	mockPicked := random.WeightedInt(mockRNG, weights)
	pseudoRNG := rand.New(rand.NewSource(324))
	pseudoPicked := random.WeightedInt(pseudoRNG, weights)
	cryptoRNG := random.CryptoRNG{}
	cryptoPicked := random.WeightedInt(cryptoRNG, weights)
	app.logger.DebugContext(
		ctx,
		fmt.Sprintf("random picked: mock: %d, pseduo: %d, crypto: %d", mockPicked, pseudoPicked, cryptoPicked),
	)

	return RandomExample{
		MockRNG:   mockPicked,
		PseudoRNG: pseudoPicked,
		CryptoRNG: cryptoPicked,
	}, nil
}
