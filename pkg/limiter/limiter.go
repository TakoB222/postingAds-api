package limiter

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	sync.RWMutex

	visitors map[string]*visitor
	limit rate.Limit
	burst int
	ttl time.Duration
}

func NewRateLimiter(rps, burst int, ttl time.Duration) *rateLimiter {
	return &rateLimiter{
		visitors: make(map[string]*visitor),
		limit: rate.Limit(rps),
		burst: burst,
		ttl: ttl,
	}
}

func (l *rateLimiter) cleanupVisitors(){

	for  {
		time.Sleep(time.Minute)

		l.RLock()
		for k, v := range l.visitors{
			if time.Since(v.lastSeen) > l.ttl {
				delete(l.visitors, k)
			}
		}
		l.Unlock()
	}

}

func (l *rateLimiter) getVisitor(ip string) *rate.Limiter{
	l.RLock()
	v, ok := l.visitors[ip]
	l.RUnlock()

	if !ok {
		limiter := rate.NewLimiter(l.limit, l.burst)
		l.RLock()
		l.visitors[ip] = &visitor{
			limiter: limiter,
			lastSeen: time.Now(),
		}
		l.Unlock()

		return limiter
	}

	v.lastSeen = time.Now()

	return v.limiter
}

func Limit(rps, burst int, ttl time.Duration) gin.HandlerFunc {
	l := NewRateLimiter(rps, burst, ttl)

	go l.cleanupVisitors()

	return func(ctx *gin.Context) {
		requestIP, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		if !l.getVisitor(requestIP).Allow() {
			ctx.AbortWithStatus(http.StatusTooManyRequests)

			return
		}

		ctx.Next()
	}
}




