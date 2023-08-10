package price_client

import "errors"

type MockPriceServiceClient struct {
	GetKeyFunc      func() string
	GetPriceFunc    func() (*PriceInfo, error)
	HealthCheckFunc func() error
}

func (m MockPriceServiceClient) GetKey() string {
	if m.GetKeyFunc != nil {
		return m.GetKeyFunc()
	}
	return "test url"
}

func (m MockPriceServiceClient) HealthCheck() error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc()
	}
	return nil
}

func (m MockPriceServiceClient) GetPrice() (*PriceInfo, error) {
	if m.GetPriceFunc != nil {
		return m.GetPriceFunc()
	}
	return &PriceInfo{}, nil
}

func NewMockPriceServiceClientSuccessCase(key string, amount int64, fraction int) *MockPriceServiceClient {
	return &MockPriceServiceClient{
		GetKeyFunc: func() string {
			return key
		},
		GetPriceFunc: func() (*PriceInfo, error) {
			return &PriceInfo{Amount: amount, Fraction: fraction}, nil
		},
	}
}

func NewMockPriceServiceClientFailCase(key string) *MockPriceServiceClient {
	return &MockPriceServiceClient{
		GetKeyFunc: func() string {
			return key
		},
		GetPriceFunc: func() (*PriceInfo, error) {
			return nil, errors.New("dummy error")
		},
	}
}
