package random

// RNG 抽象隨機數生成介面，提供整數與浮點型別的隨機數產生與洗牌功能。
type RNG interface {
	// Intn 回傳 [0, n) 區間內的均勻隨機整數
	Intn(int) int

	// Float64 回傳 [0.0, 1.0) 區間的均勻隨機浮點數
	Float64() float64

	// Shuffle 對 n 個元素進行洗牌，透過 swap(i, j) 實際交換內容
	Shuffle(n int, swap func(i, j int))
}

// Float64Range 由 rng.Float64() + 線性變換產生
func Float64Range(rng RNG, min, max float64) float64 {
	return rng.Float64()*(max-min) + min
}

// WeightedInt 適用於 int 權重的加權抽樣
func WeightedInt(rng RNG, weights []int) int {
	sum := 0
	for _, w := range weights {
		sum += w
	}
	if sum == 0 {
		return 0
	}
	threshold := rng.Intn(sum)
	accum := 0
	for i, w := range weights {
		accum += w
		if threshold < accum {
			return i
		}
	}
	return len(weights) - 1
}
