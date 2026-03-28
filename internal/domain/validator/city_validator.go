package validator

type CityValidator interface {
	Validate(city string) (*string, error)
}
