package mocks

import "github.com/stretchr/testify/mock"

type MockCityValidator struct {
	mock.Mock
}

func (m *MockCityValidator) Validate(city string) (*string, error) {
	args := m.Called(city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	validatedCity := args.Get(0).(string)
	return &validatedCity, args.Error(1)
}
