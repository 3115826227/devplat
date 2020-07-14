package docker

import (
	"devplat/src/config"
	"devplat/src/log"
	"fmt"
)

type CliManager ContainerManager

var cliManager *CliManager

func initCliManager() {
	var conatainerName = cliName
	cliManager = &CliManager{
		Env: []string{
			"GOPATH=/opt/gopath",
			"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
			"FABRIC_LOGGING_SPEC=DEBUG",
			"CORE_PEER_ID=cli",
			fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", peerName),
			"CORE_PEER_LOCALMSPID=DEFAULT",
			"CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp",
		},
		Cmd:           []string{"/bin/bash", "-c", "./script.sh"},
		Image:         "hyperledger/fabric-tools:latest",
		WorkingDir:    "/opt/gopath/src/chaincodedev",
		ContainerName: conatainerName,
		Volumes: []string{
			config.WorkPath + "/deploy/msp:/etc/hyperledger/msp",
			"/var/run:/host/var/run",
			config.WorkPath + "/deploy/chaincode:/opt/gopath/src/chaincodedev/chaincode",
			config.WorkPath + "/deploy:/opt/gopath/src/chaincodedev/",
		},
		Ports: map[string][]Port{},
	}
}

func (manager *CliManager) AddConfig() bool {
	return (*ContainerManager)(manager).AddConfig()
}

func (manager *CliManager) Run() bool {
	log.Logger.Info("cli container start to create")
	return (*ContainerManager)(manager).Run()
}

func (manager *CliManager) Restart() (ok bool) {
	return (*ContainerManager)(manager).Restart()
}

func (manager *CliManager) StopAndRemove() (ok bool) {
	log.Logger.Info("cli container start to stop and remove")
	return (*ContainerManager)(manager).StopAndRemove()
}

func (manager *CliManager) Stop() (ok bool) {
	return (*ContainerManager)(manager).Stop()
}

func (manager *CliManager) Remove() (ok bool) {
	return (*ContainerManager)(manager).Remove()
}

func (manager *CliManager) Exec(cmd []string) (content string, ok bool) {
	return (*ContainerManager)(manager).Exec(cmd)
}

func (manager *CliManager) GetStatus() ContainerManager {
	return (ContainerManager)(*manager)
}
