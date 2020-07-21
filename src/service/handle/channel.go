package handle

import (
	"devplat/src/log"
	"devplat/src/service/app"
	"github.com/gin-gonic/gin"
)

func GetChannelsHandle(c *gin.Context) {
	channels, err := app.GetChannels()
	if err != nil {
		log.Logger.Error(err.Error())
		ErrorResp(c, paramError)
		return
	}
	SuccessResp(c, "", channels)
}
