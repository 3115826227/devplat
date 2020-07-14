package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	paramError = "参数错误"
)

func SuccessResp(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"ret_code": 0, "ret_msg": message, "data": data})
}

func ErrorResp(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"ret_code": 400,
		"ret_msg":  message,
		"data":     struct{}{},
	})
}
