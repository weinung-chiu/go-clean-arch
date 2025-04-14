package random

import "testing"

// TestFloat64Range 檢查 Float64Range 計算是否正確轉換區間
func TestFloat64Range(t *testing.T) {
	rng := &MockRNG{NextFloat: 0.5}
	result := Float64Range(rng, 10.0, 2.0)
	expected := 2.0 + 0.5*(10.0-2.0) // = 6.0

	if result != expected {
		t.Errorf("Float64Range failed, expected %.2f, got %.2f", expected, result)
	}
}

// TestWeightedInt 檢查加權抽樣是否根據 threshold 落點回傳正確 index
func TestWeightedInt(t *testing.T) {
	weights := []int{2, 3, 5} // 累積為 [2, 5, 10]

	tests := []struct {
		threshold int
		expected  int
	}{
		{0, 0}, // 0 < 2
		{2, 1}, // 2 >= 2, < 5
		{9, 2}, // 9 >= 5, < 10
	}

	for _, tt := range tests {
		rng := &MockRNG{NextInt: tt.threshold}
		result := WeightedInt(rng, weights)
		if result != tt.expected {
			t.Errorf("WeightedInt(%d) failed, expected %d, got %d", tt.threshold, tt.expected, result)
		}
	}
}

// TestWeightedInt_ZeroWeight 測試全為 0 權重時是否回傳 0（邊界情境）
func TestWeightedInt_ZeroWeight(t *testing.T) {
	rng := &MockRNG{NextInt: 0}
	weights := []int{0, 0, 0}
	result := WeightedInt(rng, weights)

	if result != 0 {
		t.Errorf("WeightedInt with zero weights failed, expected 0, got %d", result)
	}
}
