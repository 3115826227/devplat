#!/bin/bash

## 设置配置环境
export FABRIC_CFG_PATH=$PWD
## 生成证书
bin/cryptogen generate --config=./crypto-config.yaml
## 生成创世块
bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
## 生成配置区块
bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
## 更新锚节点配置
bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
sleep 3
## 启动服务
docker-compose -f docker-compose/docker-compose.yaml up -d