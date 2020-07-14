package handle

import (
	"devplat/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupPlatHandle(c *gin.Context) {
	service.GetDevPlatController().Setup()
	status, containers := service.GetDevPlatController().GetContainers()
	SuccessResp(c, "", gin.H{
		"status":     status,
		"containers": containers,
	})
}

func CleanPlatHandle(c *gin.Context) {
	service.GetDevPlatController().Clean()
	GetChaincodeProvider().Clean()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "清理成功",
	})
}

func GetStatusHandle(c *gin.Context) {
	status, containers := service.GetDevPlatController().GetContainers()
	SuccessResp(c, "", gin.H{
		"status":     status,
		"containers": containers,
	})
}
