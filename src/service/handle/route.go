package handle

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	engine.GET("/setup", SetupPlatHandle)
	engine.GET("/clean", CleanPlatHandle)
	engine.GET("/status", GetStatusHandle)

	engine.POST("/chaincode/install", ChaincodeInstallHandle)
	engine.POST("/chaincode/instantiate", ChaincodeInstantiateHandle)
	engine.POST("/chaincode/invoke", ChaincodeInvokeHandle)
}
