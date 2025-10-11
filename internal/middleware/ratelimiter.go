package middleware

import (
	"time"

	"github.com/didip/tollbooth/v7"
	tollbooth_gin "github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func RateLimiter() gin.HandlerFunc {
	limiter := tollbooth.NewLimiter(5, nil) // 5 req/sec
	limiter.SetTokenBucketExpirationTTL(time.Minute)
	limiter.SetMessage("You have reached the maximum number of requests. Please try again later.")
	limiter.SetMessageContentType("application/json")

	return tollbooth_gin.LimitHandler(limiter)
}
