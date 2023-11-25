package immutable

type Float64Map interface {
	Get(key string) float64
	Has(key string) (float64, bool)
}

type float64Map struct {
	m map[string]float64
}

func (m *float64Map) Get(key string) float64 {
	return m.m[key]
}

func (m *float64Map) Has(key string) (float64, bool) {
	v, ok := m.m[key]

	return v, ok
}

func MakeImmutable(m map[string]float64) Float64Map {
	mCopy := make(map[string]float64)

	for k, v := range m {
		mCopy[k] = v
	}

	return &float64Map{m: mCopy}
}
