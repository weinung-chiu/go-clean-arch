package random

type MockRNG struct {
	NextInt   int
	NextFloat float64
}

func (m *MockRNG) Intn(n int) int {
	if m.NextInt < n {
		return m.NextInt
	}
	return n - 1 // safeguard
}

func (m *MockRNG) Float64() float64 {
	return m.NextFloat
}

func (m *MockRNG) Shuffle(n int, swap func(i, j int)) {
	panic("implement me")
}
