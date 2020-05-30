package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/health"
)

const HealthGroup = "/health"
const helloEndpoint = "/hello"

func BuildHelloEndpoint(handlers ...gin.HandlerFunc) Endpoint {
	return Endpoint{
		URL:      helloEndpoint,
		Handlers: append(handlers, hello),
	}
}

func hello(ctx *gin.Context) {
	health.Hello(ctx)
}
