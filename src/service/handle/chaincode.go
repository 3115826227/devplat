package handle

import (
	"devplat/src/config"
	"devplat/src/log"
	"devplat/src/service/app"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type InstallChaincodeReq struct {
	Result bool   `json:"result"`
	Log    string `json:"log"`
}

func InstallChaincodeHandle(c *gin.Context) {
}

type InvokeChaincodeReq struct {
	Name         string `json:"name" binding:"required"`
	FunctionName string `json:"function_name" binding:"required"`
	Args         []interface{}
}

func InvokeChaincodeHandle(c *gin.Context) {
	var req InvokeChaincodeReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}

	var args [][]byte
	for _, arg := range req.Args {
		switch arg.(type) {
		case string:
			args = append(args, []byte((arg).(string)))
		default:
			data, err := json.Marshal(arg)
			if err != nil {
				log.Logger.Error(err.Error())
				ErrorResp(c, paramError)
				return
			}
			args = append(args, data)
		}
	}
	var cfg = config.Config
	payload, err := app.InvokeCCRequest(cfg.ChannelName, req.Name, req.FunctionName, args)
	if err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, err.Error())
		return
	}
	SuccessResp(c, "", string(payload))
}
