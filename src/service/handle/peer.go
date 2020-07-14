package handle

import (
	"devplat/src/log"
	"devplat/src/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type ChaincodeInstallReq struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type ChaincodeInstallRsp struct {
	Result bool   `json:"result"`
	Log    string `json:"log"`
}

func ChaincodeInstallHandle(c *gin.Context) {
	var req ChaincodeInstallReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var cmdStr = fmt.Sprintf("peer chaincode install -p chaincodedev/chaincode/%v -n %s -v %s", req.Name, req.Name, req.Version)
	var cmd = strings.Split(cmdStr, " ")
	content, ok := service.GetDockerManager().CliManager.Exec(cmd)
	SuccessResp(c, "", ChaincodeInstallRsp{
		Result: ok,
		Log:    content,
	})
}

type ChaincodeInstantiateReq struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Args    []interface{} `json:"args"`
}

func ChaincodeInstantiateHandle(c *gin.Context) {
	var req ChaincodeInstantiateReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var argsStr string
	for _, arg := range req.Args {
		switch arg.(type) {
		case string:
			argsStr = fmt.Sprintf(`%v,"%v"`, argsStr, arg)
		case map[string]interface{}:
			data, _ := json.Marshal(arg)
			newData := strings.Replace(string(data), `"`, `\"`, len(string(data)))
			argsStr = fmt.Sprintf(`%v,"%v"`, argsStr, newData)
		default:
			ErrorResp(c, paramError)
			return
		}
	}
	var cmdStr = fmt.Sprintf(`peer chaincode instantiate -n %s -v %s -c {"Args":["init"%v]} -C myc`,
		req.Name, req.Version, argsStr)
	var cmd = strings.Split(cmdStr, " ")
	fmt.Println(cmdStr, cmd)
	content, ok := service.GetDockerManager().CliManager.Exec(cmd)
	SuccessResp(c, "", ChaincodeInstallRsp{
		Result: ok,
		Log:    content,
	})
}

type ChaincodeInvokeReq struct {
	Name         string        `json:"name" binding:"required"`
	FunctionName string        `json:"function_name" binding:"required"`
	Args         []interface{} `json:"args"`
}

type ChaincodeInvokeRsp struct {
	Result   bool   `json:"result"`
	Response string `json:"response"`
}

func ChaincodeInvokeHandle(c *gin.Context) {
	var req ChaincodeInvokeReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var argsStr string
	for _, arg := range req.Args {
		switch arg.(type) {
		case string:
			argsStr = fmt.Sprintf(`%v,"%v"`, argsStr, arg)
		case map[string]interface{}:
			data, _ := json.Marshal(arg)
			newData := strings.Replace(string(data), `"`, `\"`, len(string(data)))
			argsStr = fmt.Sprintf(`%v,"%v"`, argsStr, newData)
		default:
			ErrorResp(c, paramError)
			return
		}
	}
	var cmdStr = fmt.Sprintf(`peer chaincode invoke -n %s -c {"Args":["%v"%v]} -C myc`,
		req.Name, req.FunctionName, argsStr)
	var cmd = strings.Split(cmdStr, " ")
	fmt.Println(cmdStr, cmd)
	content, ok := service.GetDockerManager().CliManager.Exec(cmd)
	SuccessResp(c, "", ChaincodeInvokeRsp{
		Result:   ok,
		Response: content,
	})
}
