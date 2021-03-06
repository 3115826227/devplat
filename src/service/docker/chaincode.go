package docker

import (
	"devplat/src/config"
	"devplat/src/log"
	"fmt"
)

type ChaincodeManager ContainerManager

var chaincodeManager *ChaincodeManager

func initChaincodeManager() {
	var conatainerName = chaincodeName
	chaincodeManager = &ChaincodeManager{
		Env: []string{
			"GOPATH=/opt/gopath",
			"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
			"FABRIC_LOGGING_SPEC=DEBUG",
			"CORE_PEER_ID=example_cc",
			fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", peerName),
			"CORE_PEER_LOCALMSPID=DEFAULT",
			"CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp",
		},
		Cmd:           []string{"/bin/bash", "-c", "sleep 6000000"},
		Image:         "hyperledger/fabric-tools:latest",
		WorkingDir:    "/opt/gopath/src/chaincode",
		ContainerName: conatainerName,
		Volumes: []string{
			config.WorkPath + "/deploy/msp:/etc/hyperledger/msp",
			"/var/run:/host/var/run",
			config.WorkPath + "/deploy/chaincode:/opt/gopath/src/chaincode",
		},
		Ports: map[string][]Port{},
	}
}

func (manager *ChaincodeManager) AddConfig() bool {
	return (*ContainerManager)(manager).AddConfig()
}

func (manager *ChaincodeManager) Run() bool {
	log.Logger.Info("chaincode container start to create")
	return (*ContainerManager)(manager).Run()
}

func (manager *ChaincodeManager) Restart() (ok bool) {
	return (*ContainerManager)(manager).Restart()
}

func (manager *ChaincodeManager) StopAndRemove() (ok bool) {
	log.Logger.Info("chaincode container start to stop and remove")
	return (*ContainerManager)(manager).StopAndRemove()
}

func (manager *ChaincodeManager) Stop() (ok bool) {
	return (*ContainerManager)(manager).Stop()
}

func (manager *ChaincodeManager) Remove() (ok bool) {
	return (*ContainerManager)(manager).Remove()
}

func (manager *ChaincodeManager) Exec(cmd []string) (content string, ok bool) {
	return (*ContainerManager)(manager).Exec(cmd)
}

func (manager *ChaincodeManager) GetStatus() ContainerManager {
	return (ContainerManager)(*manager)
}
