package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimit struct {
	mutex   sync.Mutex
	iprates map[string]*rate.Limiter
}

var limiter *IPLimit

func init() {
	limiter = &IPLimit{
		iprates: make(map[string]*rate.Limiter),
	}
}

func IPLimitRaterMiddleware(c *gin.Context) {

	ip := c.ClientIP()
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()
	l, ok := limiter.iprates[ip]
	if !ok {
		l = rate.NewLimiter(1, 10)
		limiter.iprates[ip] = l
	}

	if !l.Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "too many requests",
		})
		return
	}
	c.Next()
}
