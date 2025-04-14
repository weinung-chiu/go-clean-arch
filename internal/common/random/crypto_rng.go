package random

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math/big"
)

var ErrRandRead = errors.New("crypto/rand: failed to read random bytes")

// CryptoRNG 是使用 crypto/rand 實作的強隨機數產生器
type CryptoRNG struct{}

// Intn 回傳 [0, n) 內的隨機整數，使用 crypto/rand 實作
func (c CryptoRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	bigN := big.NewInt(int64(n))
	result, err := rand.Int(rand.Reader, bigN)
	if err != nil {
		panic(ErrRandRead)
	}
	return int(result.Int64())
}

// Float64 回傳 [0.0, 1.0) 間的均勻分佈隨機浮點數，採用 53-bit 精度
func (c CryptoRNG) Float64() float64 {
	// 讀取 8 bytes，轉換為 uint64
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(ErrRandRead)
	}
	u := binary.BigEndian.Uint64(buf[:]) >> 11 // 保留高 53 位元
	return float64(u) / (1 << 53)
}

// Shuffle 實作與 math/rand 相同邏輯，但使用 CryptoRNG 提供的 Intn
func (c CryptoRNG) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("CryptoRNG: negative count in Shuffle")
	}
	for i := n - 1; i > 0; i-- {
		j := c.Intn(i + 1)
		swap(i, j)
	}
}
