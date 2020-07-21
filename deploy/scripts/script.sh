#!/bin/bash

set -e
#CORE_PEER_LOCALMSPID="Org1MSP"
#CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
#CORE_PEER_ADDRESS=peer0.org1.example.com:7051

peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx
res=$?
echo "===================== Channel is created successfully ===================== "
sleep 2
peer channel join -b ./mychannel.block
sleep 600000
exit 0
