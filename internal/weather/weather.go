package weather

import "context"

type IWeather interface {
	IsRainy(ctx context.Context, location string, threshold int) (bool, error)
}
