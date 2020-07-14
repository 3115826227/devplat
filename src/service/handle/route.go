package handle

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	engine.GET("/setup", SetupPlatHandle)
	engine.GET("/clean", CleanPlatHandle)
	engine.GET("/status", GetStatusHandle)

	engine.GET("/chaincode", ChaincodeGetHandle)
	engine.POST("/chaincode/install", ChaincodeInstallHandle)
	engine.POST("/chaincode/install/feedback", ChaincodeInstallFeedbackHandle)
	engine.POST("/chaincode/instantiate", ChaincodeInstantiateHandle)
	engine.POST("/chaincode/instantiate/feedback", ChaincodeInstantiateFeedbackHandle)
	engine.POST("/chaincode/invoke", ChaincodeInvokeHandle)
}
