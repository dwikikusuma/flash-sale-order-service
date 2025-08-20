package infrastructure

import (
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"time"
)

func GetRateLimiter() middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(1),
				Burst:     5,
				ExpiresIn: 1 * time.Minute,
			}),
	}
}
