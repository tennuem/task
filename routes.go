package main

import "github.com/gin-gonic/gin"

func routes(r *gin.Engine) {
	r.PUT("/add", addTask)
	r.GET("/get", getTask)
	r.DELETE("/delete", deleteTask)
	r.GET("/responses", getResponses)
}
