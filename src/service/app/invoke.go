package app

import (
	Config "devplat/src/config"
	"devplat/src/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ChaincodeInvokeProvider struct {
	Peers  []string
	Orders []string
}

var ccip *ChaincodeInvokeProvider

func init() {
	ccip = NewChaincodeInvokeProvider()
	go ccip.Start()
	cfg := Config.Config
	cli = getClient(cfg.SdkCfgPath, cfg.OrgName, cfg.UserName)
}

func NewChaincodeInvokeProvider() *ChaincodeInvokeProvider {
	return &ChaincodeInvokeProvider{}
}

func (ccip *ChaincodeInvokeProvider) Start() {
	ccip.Peers = Config.Config.Peers
	ccip.Orders = Config.Config.Orderer
}

/*
	Invoke或Query方法调用链码
*/
func invokeOrQueryCC(invoke bool, channelID, chaincode, functionName string, args [][]byte) (payload []byte, err error) {
	if len(ccip.Peers) == 0 {
		err = errors.New("no peer to invoke")
		log.Logger.Error(err.Error())
		return nil, err
	}
	req := &channel.Request{
		ChaincodeID:  chaincode,
		Fcn:          functionName,
		Args:         args,
		TransientMap: Config.TransientMap,
	}
	cfg := Config.Config
	var resp channel.Response
	resp, err = invokeChainCode(req, ccip.Peers, ccip.Orders, channelID,
		cfg.SdkCfgPath, cfg.OrgName, cfg.UserName, invoke)
	return resp.Payload, err
}

/*
	Query方式调用链码
*/
func QueryCCRequest(channelID, chaincode, functionName string, args [][]byte) (payload []byte, err error) {
	return invokeOrQueryCC(false, channelID, chaincode, functionName, args)
}

/*
	Invoke方式调用链码
*/
func InvokeCCRequest(channelID, chaincode, functionName string, args [][]byte) (payload []byte, err error) {
	return invokeOrQueryCC(true, channelID, chaincode, functionName, args)
}

var cli *client

type client struct {
	SDK            *fabsdk.FabricSDK
	ChannelClients map[string]*channel.Client
	resmgmtClient  *resmgmt.Client
	orgName        string
	userName       string
}

func GetClient() *client {
	return cli
}

func invokeChainCode(req *channel.Request, peerEndpoints, ordererEndpoint []string, channelID, sdkCfg, orgName, userName string, invoke bool) (channel.Response, error) {
	cli = getClient(sdkCfg, orgName, userName)
	if cli == nil {
		return channel.Response{}, errors.New("can not get controller client")
	}

	return cli.invokeChainCode(channelID, peerEndpoints, ordererEndpoint, req, invoke)
}

/*
	获取客户端
*/
func getClient(sdkCfg, orgName, userName string) *client {
	if cli != nil {
		return cli
	}

	var (
		err           error
		sdk           *fabsdk.FabricSDK
		resmgmtClient *resmgmt.Client
	)

	if sdk, err = fabsdk.New(config.FromFile(sdkCfg)); err != nil {
		log.Logger.Error("create fabric sdk from file failed",
			zap.String("config file", sdkCfg), zap.Error(err))
		return nil
	}

	rcp := sdk.Context(fabsdk.WithOrg(orgName), fabsdk.WithUser(userName))
	if resmgmtClient, err = resmgmt.New(rcp); err != nil {
		log.Logger.Error("create resource management client failed",
			zap.String("org", orgName), zap.String("user", userName),
			zap.Error(err))
		return nil
	}

	cli = &client{
		SDK:            sdk,
		ChannelClients: make(map[string]*channel.Client),
		resmgmtClient:  resmgmtClient,
		orgName:        orgName,
		userName:       userName,
	}
	return cli
}

/*
	获取具体的channel客户端
*/
func (c *client) getChannelClient(channelID string) (*channel.Client, error) {
	if cc, ok := c.ChannelClients[channelID]; ok {
		return cc, nil
	}
	var (
		err error
		cc  *channel.Client
	)

	ccp := c.SDK.ChannelContext(channelID, fabsdk.WithUser(c.userName))
	if cc, err = channel.New(ccp); err != nil {
		return nil, errors.WithMessage(err, "failed to create channel client")
	}
	c.ChannelClients[channelID] = cc

	return cc, nil
}

/*
	调用链码
*/
func (c *client) invokeChainCode(channelID string, peerEndpoints, ordererEndpoints []string, req *channel.Request, invoke bool) (channel.Response, error) {
	var (
		ccli *channel.Client
		err  error
	)

	if ccli, err = c.getChannelClient(channelID); err != nil {
		return channel.Response{}, errors.WithMessage(err, "invoke chaincode")
	}
	var opts []channel.RequestOption
	if peerEndpoints != nil {
		opts = append(opts, channel.WithTargetEndpoints(peerEndpoints...))
	}
	//for _, orderer := range ordererEndpoints {
	//	opts = append(opts, channel.WithOrdererEndpoint(orderer))
	//}
	if invoke {
		return ccli.Execute(*req, opts...)
	}
	return ccli.Query(*req, opts...)
}
