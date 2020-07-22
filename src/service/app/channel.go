package app

import (
	"devplat/src/config"
	"devplat/src/log"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

func GetChannels() (channels []string, err error) {
	var response *peer.ChannelQueryResponse
	var opts []resmgmt.RequestOption
	opts = append(opts, resmgmt.WithTargetEndpoints(config.Config.Peers[0]))
	response, err = GetClient().resmgmtClient.QueryChannels(opts...)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, c := range response.Channels {
		channels = append(channels, c.ChannelId)
	}
	return
}
