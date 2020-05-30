package endpoints

import "github.com/gin-gonic/gin"

type Endpoint struct{
	URL string
	Handlers []gin.HandlerFunc
}

