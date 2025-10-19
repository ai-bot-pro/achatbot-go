package middleware

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter Token Bucket 速率限制器
type RateLimiter struct {
	enabled              bool
	limiters             map[string]*rate.Limiter
	mu                   sync.RWMutex
	r                    rate.Limit // 每秒向桶中放多少token (每秒最大请求数目)
	b                    int        // token桶容量大小
	maxConns             int
	connCount            int32
	cleanupIntervalTimeS int // 清理间隔时间,检查是否token桶是满的，满则有段时间未用，可删除释放对应ip limiter
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(
	enabled bool,
	requestsPerSecond, burstSize, maxConnections int,
) *RateLimiter {
	return &RateLimiter{
		enabled:              enabled,
		limiters:             make(map[string]*rate.Limiter),
		r:                    rate.Limit(requestsPerSecond),
		b:                    burstSize,
		maxConns:             maxConnections,
		cleanupIntervalTimeS: 3,
	}
}

func NewDefaultRateLimiter() *RateLimiter {
	return &RateLimiter{
		enabled:              true,
		limiters:             make(map[string]*rate.Limiter),
		r:                    rate.Limit(1024),
		b:                    1024,
		maxConns:             1024,
		cleanupIntervalTimeS: 3,
	}
}

func (rl *RateLimiter) WithEnable(enabled bool) *RateLimiter {
	rl.enabled = enabled
	return rl
}

func (rl *RateLimiter) WithRequestsPerSecond(requestsPerSecond int) *RateLimiter {
	rl.r = rate.Limit(requestsPerSecond)
	return rl
}

func (rl *RateLimiter) WithBurstSize(burstSize int) *RateLimiter {
	rl.b = burstSize
	return rl
}

func (rl *RateLimiter) WithMaxConns(maxConns int) *RateLimiter {
	rl.maxConns = maxConns
	return rl
}

func (rl *RateLimiter) WithCleanupIntervalTimeS(cleanupIntervalTimeS int) *RateLimiter {
	rl.cleanupIntervalTimeS = cleanupIntervalTimeS
	return rl
}

// getLimiter 获取或创建IP的限制器
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()
	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check in case another goroutine created it while we were waiting for the lock.
	if limiter, exists = rl.limiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rl.r, rl.b)
	rl.limiters[ip] = limiter
	return limiter
}

// cleanupLimiters 清理过期的限制器
func (rl *RateLimiter) cleanupLimiters() {
	if rl.cleanupIntervalTimeS < 3 {
		rl.cleanupIntervalTimeS = 3
	}
	ticker := time.NewTicker(time.Duration(rl.cleanupIntervalTimeS) * time.Second)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			for ip, limiter := range rl.limiters {
				//allow := limiter.AllowN(time.Now(), rl.b)
				availableTokenNum := limiter.Tokens()
				allow := availableTokenNum >= float64(rl.b) // for float compare
				if allow {
					// 检查是否token桶是满的，满则有段时间未用，可删除释放对应ip limiter
					delete(rl.limiters, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}

// Middleware 速率限制中间件
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	// 如果限速器未启用，直接跳过
	if !rl.enabled {
		return next
	}

	// 启动清理协程
	rl.cleanupLimiters()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查连接数限制
		currentConns := atomic.AddInt32(&rl.connCount, 1)
		if currentConns > int32(rl.maxConns) {
			atomic.AddInt32(&rl.connCount, -1) // Decrement back as we are rejecting this connection.
			http.Error(w, "Too many connections", http.StatusTooManyRequests)
			return
		}

		// 连接结束时减少计数
		defer atomic.AddInt32(&rl.connCount, -1)

		// 获取客户端IP
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			// X-Forwarded-For can be a comma-separated list of IPs. The first one is the original client.
			parts := strings.Split(forwarded, ",")
			ip = strings.TrimSpace(parts[0])
		}

		// 检查客户端IP请求速率限制
		limiter := rl.getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) GetClientIPs() (keys []string) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	for key := range rl.limiters {
		keys = append(keys, key)
	}
	return
}

// GetStats 获取统计信息
func (rl *RateLimiter) GetStats() map[string]any {
	currentConns := atomic.LoadInt32(&rl.connCount)

	rl.mu.RLock()
	activeLimiters := len(rl.limiters)
	rl.mu.RUnlock()

	return map[string]any{
		"enabled":             rl.enabled,
		"active_limiters":     activeLimiters,
		"current_connections": currentConns,
		"max_connections":     rl.maxConns,
		"requests_per_second": float64(rl.r),
		"burst_size":          rl.b,
	}
}
