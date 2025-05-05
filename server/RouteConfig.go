package server

import "github.com/gin-gonic/gin"

type RouteConfig struct {
	Path       string
	Method     string
	HandleFunc func(c *gin.Context)
	Middleware []gin.HandlerFunc
}
