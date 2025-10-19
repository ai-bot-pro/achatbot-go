package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(true, 10, 5, 100)
	assert.NotNil(t, rl)
	assert.True(t, rl.enabled)
	assert.Equal(t, rate.Limit(10), rl.r)
	assert.Equal(t, 5, rl.b)
	assert.Equal(t, 100, rl.maxConns)
	assert.NotNil(t, rl.limiters)
}

func TestGetLimiter(t *testing.T) {
	rl := NewRateLimiter(true, 10, 5, 100)
	ip := "127.0.0.1"

	// 第一次获取
	limiter1 := rl.getLimiter(ip)
	assert.NotNil(t, limiter1)

	// 第二次获取同一个IP的限制器
	limiter2 := rl.getLimiter(ip)
	assert.NotNil(t, limiter2)

	// 应该是同一个实例
	assert.Equal(t, limiter1, limiter2)
}

func TestGetStats(t *testing.T) {
	rl := NewRateLimiter(true, 10, 5, 100)

	stats := rl.GetStats()
	assert.NotNil(t, stats)
	assert.True(t, stats["enabled"].(bool))
	assert.Equal(t, 0, stats["active_limiters"])
	assert.Equal(t, int32(0), stats["current_connections"])
	assert.Equal(t, 100, stats["max_connections"])
	assert.Equal(t, float64(10), stats["requests_per_second"])
	assert.Equal(t, 5, stats["burst_size"])
}

func TestMiddlewareDisabled(t *testing.T) {
	rl := NewRateLimiter(false, 1, 1, 1) // 禁用限流器

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Body.String())
}

func TestMiddlewareConnectionLimit(t *testing.T) {
	rl := NewRateLimiter(true, 10, 5, 2) // 最大连接数为2

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // 模拟处理时间
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	// 创建多个并发请求
	var wg sync.WaitGroup
	results := make(chan int, 5)

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			results <- resp.Code
		}()
	}

	wg.Wait()
	close(results)

	// 统计结果
	tooManyRequests := 0
	ok := 0

	for code := range results {
		switch code {
		case http.StatusTooManyRequests:
			tooManyRequests++
		case http.StatusOK:
			ok++
		}
	}

	// 应该有2个成功，3个被限制
	assert.Equal(t, 2, ok)
	assert.Equal(t, 3, tooManyRequests)
}

func TestMiddlewareRateLimit(t *testing.T) {
	rl := NewRateLimiter(true, 1, 1, 100) // 1请求/秒，突发1个

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1"

	// 第一个请求应该成功
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req)
	assert.Equal(t, http.StatusOK, resp1.Code)

	// 紧接着的第二个请求应该被限制
	resp2 := httptest.NewRecorder()
	handler.ServeHTTP(resp2, req)
	assert.Equal(t, http.StatusTooManyRequests, resp2.Code)
}

func TestMiddlewareWithXForwardedFor(t *testing.T) {
	rl := NewRateLimiter(true, 1, 1, 100)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	clientIP := "192.168.1.1"
	req.Header.Set("X-Forwarded-For", clientIP)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	clientIps := rl.GetClientIPs()
	assert.Equal(t, 1, len(clientIps))
	assert.Equal(t, clientIP, clientIps[0])
}

// 测试 cleanupLimiters 函数
func TestMiddlewareWithCleanupIntervalTimeS(t *testing.T) {
	cleanupIntervalTimeS := 3
	rl := NewRateLimiter(true, 1, 1, 100).WithCleanupIntervalTimeS(cleanupIntervalTimeS)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	clientIP := "192.168.1.1"
	req.Header.Set("X-Forwarded-For", clientIP)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	clientIps := rl.GetClientIPs()
	assert.Equal(t, 1, len(clientIps))
	assert.Equal(t, clientIP, clientIps[0])

	time.Sleep((time.Duration(cleanupIntervalTimeS) + 1) * time.Second)
	assert.Equal(t, 0, len(rl.GetClientIPs()))
}

func TestMultiCleanupLimiters(t *testing.T) {
	// 创建一个速率限制器，设置较小的限制以便测试, 每秒放一个token,token桶容量为1, 最大并发3个请求
	cleanupIntervalTimeS := 3
	rl := NewRateLimiter(true, 1, 1, 3).WithCleanupIntervalTimeS(cleanupIntervalTimeS)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	// 添加几个限制器
	ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4"}
	for _, ip := range ips {
		for i := range 2 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Forwarded-For", ip)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if i >= 1 {
				assert.Equal(t, http.StatusTooManyRequests, resp.Code)
			} else {
				assert.Equal(t, http.StatusOK, resp.Code)
			}
		}
	}

	// 检查限制器是否已添加
	countBefore := len(rl.GetClientIPs())
	assert.Equal(t, len(ips), countBefore)
	stats := rl.GetStats()
	assert.Equal(t, len(ips), stats["active_limiters"])

	// check cleeanup
	time.Sleep((time.Duration(cleanupIntervalTimeS) + 1) * time.Second)
	assert.Equal(t, 0, len(rl.GetClientIPs()))
}

func TestMultiConcurrentCleanupLimiters(t *testing.T) {
	// 创建一个速率限制器，设置较小的限制以便测试, 每秒放一个token,token桶容量为1, 最大并发3个请求
	cleanupIntervalTimeS := 3
	maxConns := 3
	rl := NewRateLimiter(true, 1, 1, maxConns).WithCleanupIntervalTimeS(cleanupIntervalTimeS)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//slow process (e.g.: slow sql/kv(RAG) io or local asr/llm/tts DNN model inference)
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := rl.Middleware(nextHandler)

	// 添加几个限制器
	var wg sync.WaitGroup
	ips := []string{}
	for i := range maxConns {
		ips = append(ips, fmt.Sprintf("192.168.1.%d", i))
	}
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			for i := range 2 {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Forwarded-For", ip)
				resp := httptest.NewRecorder()
				handler.ServeHTTP(resp, req)
				if i >= 1 {
					assert.Equal(t, http.StatusTooManyRequests, resp.Code)
				} else {
					assert.Equal(t, http.StatusOK, resp.Code)
				}
			}
		}(ip)
	}
	time.Sleep(10 * time.Millisecond) // wait ips request goroutine to schedule execute

	otherFakeIPs := []string{}
	for i := range 100 {
		otherFakeIPs = append(otherFakeIPs, fmt.Sprintf("192.168.0.%d", i))
	}
	for _, ip := range otherFakeIPs {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			for range 2 {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Forwarded-For", ip)
				resp := httptest.NewRecorder()
				handler.ServeHTTP(resp, req)
				assert.Equal(t, http.StatusTooManyRequests, resp.Code)
			}
		}(ip)
	}
	wg.Wait()

	// 检查限制器是否已添加
	countBefore := len(rl.GetClientIPs())
	assert.Equal(t, len(ips), countBefore)
	stats := rl.GetStats()
	assert.Equal(t, len(ips), stats["active_limiters"])

	// check cleeanup
	time.Sleep((time.Duration(cleanupIntervalTimeS) + 1) * time.Second)
	assert.Equal(t, 0, len(rl.GetClientIPs()))
}
