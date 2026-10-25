package engine

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// RateLimiter implements a smooth, thread-safe Token Bucket bandwidth limiter.
// A limit <= 0 means unlimited speed.
type RateLimiter struct {
	limitBytesPerSec atomic.Int64
	mu               sync.Mutex
	tokens           float64
	lastRefill       time.Time
}

// NewRateLimiter creates a bandwidth limiter with the specified bytes/sec cap (0 = unlimited).
func NewRateLimiter(bytesPerSec int64) *RateLimiter {
	rl := &RateLimiter{
		tokens:     float64(bytesPerSec),
		lastRefill: time.Now(),
	}
	rl.limitBytesPerSec.Store(bytesPerSec)
	return rl
}

// SetLimit updates the bandwidth cap in bytes/sec at runtime without dropping active streams.
func (rl *RateLimiter) SetLimit(bytesPerSec int64) {
	if bytesPerSec < 0 {
		bytesPerSec = 0
	}
	rl.limitBytesPerSec.Store(bytesPerSec)
	rl.mu.Lock()
	if rl.tokens > float64(bytesPerSec) {
		rl.tokens = float64(bytesPerSec)
	}
	rl.lastRefill = time.Now()
	rl.mu.Unlock()
}

// GetLimit returns the current rate limit in bytes/sec (0 = unlimited).
func (rl *RateLimiter) GetLimit() int64 {
	return rl.limitBytesPerSec.Load()
}

// Wait blocks until n bytes are permitted by the token bucket or ctx is cancelled.
func (rl *RateLimiter) Wait(ctx context.Context, n int) error {
	limit := rl.limitBytesPerSec.Load()
	if limit <= 0 || n <= 0 {
		return nil
	}

	needed := float64(n)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		rl.mu.Lock()
		limitNow := rl.limitBytesPerSec.Load()
		if limitNow <= 0 {
			rl.mu.Unlock()
			return nil
		}

		now := time.Now()
		elapsed := now.Sub(rl.lastRefill).Seconds()
		if elapsed > 0 {
			rl.tokens += elapsed * float64(limitNow)
			maxBurst := float64(limitNow) * 0.5
			if maxBurst < float64(BufferSize*2) {
				maxBurst = float64(BufferSize * 2)
			}
			if rl.tokens > maxBurst {
				rl.tokens = maxBurst
			}
			rl.lastRefill = now
		}

		if rl.tokens >= needed {
			rl.tokens -= needed
			rl.mu.Unlock()
			return nil
		}

		deficit := needed - rl.tokens
		rl.tokens = 0
		waitDuration := time.Duration((deficit / float64(limitNow)) * float64(time.Second))
		rl.mu.Unlock()

		if waitDuration < 2*time.Millisecond {
			waitDuration = 2 * time.Millisecond
		}
		if waitDuration > 250*time.Millisecond {
			waitDuration = 250 * time.Millisecond
		}

		timer := time.NewTimer(waitDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	}
}
