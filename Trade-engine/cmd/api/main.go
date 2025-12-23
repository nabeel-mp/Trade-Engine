package main

import (
	"Trade-engine/Internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.Register(r)
	r.Run(":8080")
}
