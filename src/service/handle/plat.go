package handle

import (
	"devplat/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupPlatHandle(c *gin.Context) {
	service.GetDevPlatController().Setup()
	SuccessResp(c, "", service.GetDevPlatController().GetContainers())
}

func CleanPlatHandle(c *gin.Context) {
	service.GetDevPlatController().Clean()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "清理成功",
	})
}

func GetStatusHandle(c *gin.Context) {
	SuccessResp(c, "", service.GetDevPlatController().GetContainers())
}
