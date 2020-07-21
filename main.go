package main

import (
	"context"
	"devplat/src/config"
	"devplat/src/log"
	"devplat/src/service"
	"devplat/src/service/app"
	"devplat/src/service/handle"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	//获取工作目录
	config.WorkPath, _ = os.Getwd()
	config.GoPath = os.Getenv("GOPATH")
	//初始化Controller
	service.InitDevPlatController()

	app.InitChaincodeProvider()
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Error("listen: " + err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		//清理环境
		service.GetDevPlatController().Clean()
	}

	log.Logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Error("Server Shutdown:" + err.Error())
		return
	}
	log.Logger.Info("Server exiting")
}
