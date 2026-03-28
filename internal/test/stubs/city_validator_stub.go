package stubs

type CityValidatorStub struct {
	ValidateFn func(city string) (*string, error)
}

func NewCityValidatorStub() *CityValidatorStub {
	return &CityValidatorStub{
		ValidateFn: nil,
	}
}

func (s *CityValidatorStub) Validate(city string) (*string, error) {
	if s.ValidateFn != nil {
		return s.ValidateFn(city)
	}
	return &city, nil
}
