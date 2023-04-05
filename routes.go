package main

import (
	handler "hunt/internal/handler/http"
	"hunt/internal/repository/mysql"

	"github.com/gin-gonic/gin"
)

func setupRoutes() {
	r := gin.Default()

	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}

	h := handler.New(repo)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/createUser", h.CreateUser)
		v1.POST("/startMatch", h.StartMatch)
		v1.GET("/monitor", h.QueueMonitor)
		v1.POST("/login", h.Login)
	}

	r.Run()
}
