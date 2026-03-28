package models

type Weather struct {
	Temperature float64
	Humidity    float64
	Description string
}

type WeatherDaily struct {
	Location   string
	Date       string
	MaxTempC   float64
	MinTempC   float64
	AvgTempC   float64
	WillItRain bool
	ChanceRain int
	WillItSnow bool
	ChanceSnow int
	Condition  string
	Icon       string
}

type WeatherHourly struct {
	Location   string
	Time       string
	TempC      float64
	WillItRain bool
	ChanceRain int
	WillItSnow bool
	ChanceSnow int
	Condition  string
	Icon       string
}
