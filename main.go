package main

import (
	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
)

const port = ":3000"

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/", run)
	ezap.Info("listening to ", port)
	r.Run(port)
}
