package handle

import (
	"devplat/src/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

func ArgsHandle(args []interface{}) (argsStr string, err error) {
	for _, arg := range args {
		switch arg.(type) {
		case string:
			var str = arg.(string)
			str = strings.Replace(str, " ", "", -1)
			str = strings.Replace(str, "\n", "", -1)
			str = strings.Replace(str, "\"", "\\\"", -1)
			argsStr = fmt.Sprintf(`%v,"%v"`, argsStr, str)
		case map[string]interface{}:
			data, _ := json.Marshal(arg)
			newData := strings.Replace(string(data), `"`, `\"`, len(string(data)))
			argsStr = fmt.Sprintf(`%v,"%v"`, argsStr, newData)
		default:
			err = errors.New("参数类型错误")
			log.Logger.Error(err.Error())
			return
		}
	}
	return
}
