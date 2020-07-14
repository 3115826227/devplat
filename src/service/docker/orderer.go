package docker

import (
	"devplat/src/config"
	"devplat/src/log"
	"fmt"
)

type OrdererManager ContainerManager

var ordererManager *OrdererManager

func initOrdererManager() {
	var conatainerName = ordererName
	ordererManager = &OrdererManager{
		Env: []string{
			"FABRIC_LOGGING_SPEC=debug",
			fmt.Sprintf("ORDERER_GENERAL_LISTENADDRESS=%s", conatainerName),
			"ORDERER_GENERAL_GENESISMETHOD=file",
			"ORDERER_GENERAL_GENESISFILE=orderer.block",
			"ORDERER_GENERAL_LOCALMSPID=DEFAULT",
			"ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/msp",
			"GRPC_TRACE=all=true,",
			"GRPC_VERBOSITY=debug",
		},
		Cmd:           []string{"orderer"},
		WorkingDir:    "/opt/gopath/src/github.com/hyperledger/fabric",
		ContainerName: conatainerName,
		Image:         "hyperledger/fabric-orderer:latest",
		Volumes: []string{
			config.WorkPath + "/deploy/msp:/etc/hyperledger/msp",
			config.WorkPath + "/deploy/orderer.block:/etc/hyperledger/fabric/orderer.block",
		},
		Ports: map[string][]Port{
			"7050": {
				{
					IP:   "0.0.0.0",
					Port: "7050",
				},
			},
		},
	}
}

func (manager *OrdererManager) AddConfig() bool {
	return (*ContainerManager)(manager).AddConfig()
}

func (manager *OrdererManager) Run() bool {
	log.Logger.Info("orderer container start to create")
	return (*ContainerManager)(manager).Run()
}

func (manager *OrdererManager) Restart() (ok bool) {
	return (*ContainerManager)(manager).Restart()
}

func (manager *OrdererManager) StopAndRemove() (ok bool) {
	log.Logger.Info("orderer container start to stop and remove")
	return (*ContainerManager)(manager).StopAndRemove()
}

func (manager *OrdererManager) Stop() (ok bool) {
	return (*ContainerManager)(manager).Stop()
}

func (manager *OrdererManager) Remove() (ok bool) {
	return (*ContainerManager)(manager).Remove()
}

func (manager *OrdererManager) Exec(cmd []string) (content string, ok bool) {
	return (*ContainerManager)(manager).Exec(cmd)
}

func (manager *OrdererManager) GetStatus() ContainerManager {
	return (ContainerManager)(*manager)
}
