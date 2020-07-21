package app

import (
	"devplat/src/log"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func GetChannels() (channels []string, err error) {
	var response *peer.ChannelQueryResponse
	response, err = GetClient().resmgmtClient.QueryChannels()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, c := range response.Channels {
		channels = append(channels, c.ChannelId)
	}
	return
}
