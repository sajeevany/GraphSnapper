package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sirupsen/logrus"
)

//SetCtxLogger - Sets a logger access defined by LoggerKey
func SetCtxLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(logging.LoggerKey, logger)
		c.Next()
	}
}
