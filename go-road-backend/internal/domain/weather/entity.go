package weather

import "time"

type WeatherData struct {
	Lat            float64   `json:"lat"`
	Lng            float64   `json:"lng"`
	Temperature    float64   `json:"temperature"`
	FeelsLike      float64   `json:"feels_like"`
	Humidity       int       `json:"humidity"`
	Pressure       int       `json:"pressure"`
	WeatherMain    string    `json:"weather_main"`
	WeatherDesc    string    `json:"weather_desc"`
	Icon           string    `json:"icon"`
	WindSpeed      float64   `json:"wind_speed"`
	WindDeg        int       `json:"wind_deg"`
	Visibility     int       `json:"visibility"`
	RainMm         float64   `json:"rain_mm,omitempty"`
	Clouds         int       `json:"clouds"`
	FetchedAt      time.Time `json:"fetched_at"`
}

type WeatherForecast struct {
	Lat       float64         `json:"lat"`
	Lng       float64         `json:"lng"`
	Forecasts []WeatherData   `json:"forecasts"`
	FetchedAt time.Time       `json:"fetched_at"`
}

type WeatherAlert struct {
	ID          string    `json:"id"`
	Event       string    `json:"event"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Region      string    `json:"region"`
}
