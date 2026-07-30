package weather

import "context"

type Service interface {
	GetCurrentWeather(ctx context.Context, lat, lng float64) (*WeatherData, error)
	GetForecast(ctx context.Context, lat, lng float64) (*WeatherForecast, error)
	GetAlerts(ctx context.Context, lat, lng float64) ([]WeatherAlert, error)
	GetRouteWeather(ctx context.Context, routeID string) ([]WeatherData, error)
}
