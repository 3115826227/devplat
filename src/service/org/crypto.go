package org

import (
	"devplat/src/config"
	"devplat/src/log"
	"devplat/src/utils"
	"fmt"
)

const (
	CryptoConfigFileName = "crypto-config.yaml"
	CryptoCommandName    = "./bin/cryptogen"
)

func NewOrgCrypto(org Org) {
	ordererOrgs := make([]config.OrdererOrgs, 0)
	peerOrgs := make([]config.PeerOrgs, 0)
	var cryptoConfig = config.CryptoConfig{
		OrdererOrgs: ordererOrgs,
		PeerOrgs:    peerOrgs,
	}
	fileName := fmt.Sprintf("%v-%v", org.OrgName, CryptoConfigFileName)
	config.GenerateConfig(fileName, cryptoConfig)
	if err := utils.CommonRun(CryptoCommandName, []string{"generate", "--config=./deploy/%v", fileName}); err != nil {
		log.Logger.Error(err.Error())
	}
}
