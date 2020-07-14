package main

import (
	"devplat/src/config"
	"devplat/src/service"
	"devplat/src/service/handle"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func init() {
	//获取工作目录
	config.WorkPath, _ = os.Getwd()
	config.GoPath = os.Getenv("GOPATH")
	//初始化Controller
	service.InitDevPlatController()
}

func IndexHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	engine := gin.Default()

	engine.Static("static", "static")
	engine.LoadHTMLGlob("views/*")
	engine.GET("/", IndexHandle)

	handle.Router(engine)

	engine.Run(":8080")
}
