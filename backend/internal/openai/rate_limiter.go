package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redis     *redis.Client
	key       string
	maxCalls  int
	duration  time.Duration
}

func NewRateLimiter(redis *redis.Client, maxCalls int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:     redis,
		key:       "openai:ratelimit",
		maxCalls:  maxCalls,
		duration:  duration,
	}
}

func (r *RateLimiter) Allow(ctx context.Context) error {
	pipe := r.redis.Pipeline()
	now := time.Now().UnixNano()
	windowStart := now - r.duration.Nanoseconds()

	// Remove old entries
	pipe.ZRemRangeByScore(ctx, r.key, "0", fmt.Sprintf("%d", windowStart))

	// Count requests in current window
	countCmd := pipe.ZCount(ctx, r.key, fmt.Sprintf("%d", windowStart), fmt.Sprintf("%d", now))

	// Add current request
	pipe.ZAdd(ctx, r.key, redis.Z{Score: float64(now), Member: now})

	// Set expiration
	pipe.Expire(ctx, r.key, r.duration)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %v", err)
	}

	count, err := countCmd.Result()
	if err != nil {
		return fmt.Errorf("failed to get request count: %v", err)
	}

	if count >= int64(r.maxCalls) {
		return fmt.Errorf("rate limit exceeded")
	}

	return nil
} 