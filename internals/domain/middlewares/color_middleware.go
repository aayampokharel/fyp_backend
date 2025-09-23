package middleware

import (
	"fmt"
	logger "project/package/utils/pkg"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ColorStatusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		status := c.Writer.Status()

		var color string
		switch {
		case status >= 200 && status < 300:
			color = "\033[42m" // Green
		case status >= 400 && status < 500:
			color = "\033[43m" // Yellow
		case status >= 500:
			color = "\033[41m" // Red
		default:
			color = "\033[47m" // White
		}

		reset := "\033[0m"

		// Build a **plain string** log line with colored status for console only
		logMsg := fmt.Sprintf(
			"%s %s | %s | %s %s",
			c.Request.Method,
			c.Request.URL.Path,
			color+fmt.Sprintf("  %d  ", status)+reset,
			latency,
			":"+strings.Split(c.Request.Host, ":")[1],
		)
		logger.Logger.Infoln(logMsg)
	}
}
