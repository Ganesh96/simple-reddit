package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/gin-gonic/gin"
)

type rateEntry struct {
	count       int
	windowStart time.Time
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		c.Next()
	}
}

func BodySizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	clients := map[string]rateEntry{}
	lastCleanup := time.Now()

	return func(c *gin.Context) {
		now := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		key := c.ClientIP() + "|" + c.Request.Method + "|" + path

		mu.Lock()
		if now.Sub(lastCleanup) > window {
			for clientKey, entry := range clients {
				if now.Sub(entry.windowStart) > window {
					delete(clients, clientKey)
				}
			}
			lastCleanup = now
		}

		entry := clients[key]
		if entry.windowStart.IsZero() || now.Sub(entry.windowStart) > window {
			entry = rateEntry{count: 0, windowStart: now}
		}
		entry.count++
		clients[key] = entry
		remaining := window - now.Sub(entry.windowStart)
		limited := entry.count > maxRequests
		mu.Unlock()

		if limited {
			if remaining < 0 {
				remaining = window
			}
			c.Header("Retry-After", strconv.Itoa(int(remaining.Seconds())+1))
			common.RespondWithJSON(c, http.StatusTooManyRequests, common.FORBIDDEN, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
