package handle

import (
	"devplat/src/log"
	"devplat/src/service"
	"devplat/src/service/app"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

//获取未安装链码项目
func ChaincodeUninstallGetHandle(c *gin.Context) {
	SuccessResp(c, "", app.GetChaincodeProvider().GetUninstallChaincode())
}

//获取已安装未实例化链码
func ChaincodeInstalledGetHandle(c *gin.Context) {
	SuccessResp(c, "", app.GetChaincodeProvider().GetInstalledChaincodes())
}

//获取已安装链码
func ChaincodeInstantiatedGetHandle(c *gin.Context) {
	SuccessResp(c, "", app.GetChaincodeProvider().GetInstantiatedChaincode())
}

type ChaincodeReq struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type ChaincodeInstallReq ChaincodeReq

type ChaincodeInstallRsp struct {
	Result bool   `json:"result"`
	Log    string `json:"log"`
}

//链码安装
func ChaincodeInstallHandle(c *gin.Context) {
	var req ChaincodeInstallReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var cmdStr = fmt.Sprintf("peer chaincode install -p %v -n %s -v %s", req.Name, req.Name, req.Version)
	var cmd = strings.Split(cmdStr, " ")
	content, ok := service.GetDockerManager().CliManager.Exec(cmd)
	SuccessResp(c, "", ChaincodeInstallRsp{
		Result: ok,
		Log:    content,
	})
}

type ChaincodeInstallFeedbackReq ChaincodeReq

// 链码安装反馈
func ChaincodeInstallFeedbackHandle(c *gin.Context) {
	var req ChaincodeInstallFeedbackReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var chaincodeInfo = app.ChaincodeInfo{
		Name:    req.Name,
		Version: req.Version,
	}
	app.GetChaincodeProvider().ChaincodeInstallFeedback(chaincodeInfo)
	SuccessResp(c, "", chaincodeInfo)
}

type ChaincodeInstantiateReq struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Args    []interface{} `json:"args"`
}

// 链码实例化
func ChaincodeInstantiateHandle(c *gin.Context) {
	var req ChaincodeInstantiateReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var chaincodeInfo = app.ChaincodeInfo{
		Name:    req.Name,
		Version: req.Version,
	}
	if !app.GetChaincodeProvider().JudgeChaincodeInstalled(chaincodeInfo) {
		ErrorResp(c, paramError)
		return
	}
	argsStr, err := ArgsHandle(req.Args)
	if err != nil {
		ErrorResp(c, paramError)
		return
	}
	var cmdStr = fmt.Sprintf(`peer chaincode instantiate -n %s -v %s -c {"Args":["init"%v]} -C myc`,
		req.Name, req.Version, argsStr)
	var cmd = strings.Split(cmdStr, " ")
	content, ok := service.GetDockerManager().CliManager.Exec(cmd)
	SuccessResp(c, "", ChaincodeInstallRsp{
		Result: ok,
		Log:    content,
	})
}

type ChaincodeInstantiateFeedbackReq ChaincodeInstallReq

// 链码实例化反馈
func ChaincodeInstantiateFeedbackHandle(c *gin.Context) {
	var req ChaincodeInstantiateFeedbackReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	var chaincodeInfo = app.ChaincodeInfo{
		Name:    req.Name,
		Version: req.Version,
	}
	app.GetChaincodeProvider().ChaincodeInstantiateFeedback(chaincodeInfo)
	SuccessResp(c, "", chaincodeInfo)
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

// 链码调用
func ChaincodeInvokeHandle(c *gin.Context) {
	var req ChaincodeInvokeReq
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	if !app.GetChaincodeProvider().JudgeChaincodeInstantiate(req.Name) {
		ErrorResp(c, paramError)
		return
	}
	argsStr, err := ArgsHandle(req.Args)
	if err != nil {
		ErrorResp(c, paramError)
		return
	}
	var cmdStr = fmt.Sprintf(`peer chaincode invoke -n %s -c {"Args":["%v"%v]} -C myc`,
		req.Name, req.FunctionName, argsStr)
	var cmd = strings.Split(cmdStr, " ")
	content, ok := service.GetDockerManager().CliManager.Exec(cmd)
	SuccessResp(c, "", ChaincodeInvokeRsp{
		Result:   ok,
		Response: content,
	})
}
