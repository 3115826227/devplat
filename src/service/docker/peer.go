package docker

import (
	"devplat/src/config"
	"devplat/src/log"
	"fmt"
)

type PeerManager ContainerManager

var peerManager *PeerManager

func initPeerManager() {
	var conatainerName = peerName
	peerManager = &PeerManager{
		Env: []string{
			"CORE_LEDGER_STATE_STATEDATABASE=CouchDB",
			fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=%s:5984", couchdbName),
			fmt.Sprintf("CORE_PEER_ID=%s", conatainerName),
			fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", conatainerName),
			fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:7051", conatainerName),
			"CORE_PEER_LOCALMSPID=DEFAULT",
			"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
			"FABRIC_LOGGING_SPEC=DEBUG",
			"CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp",
		},
		Cmd:           []string{"peer", "node", "start", "--peer-chaincodedev=true"},
		Image:         "hyperledger/fabric-peer:latest",
		WorkingDir:    "/opt/gopath/src/github.com/hyperledger/fabric/peer",
		ContainerName: conatainerName,
		Volumes: []string{
			config.WorkPath + "/deploy/msp:/etc/hyperledger/msp",
			"/var/run:/host/var/run",
		},
		Ports: map[string][]Port{
			"7051": {
				{
					IP:   "0.0.0.0",
					Port: "7051",
				},
			},
			"7052": {
				{
					IP:   "0.0.0.0",
					Port: "7052",
				},
			},
			"7053": {
				{
					IP:   "0.0.0.0",
					Port: "7053",
				},
			},
		},
	}
}

func (manager *PeerManager) AddConfig() bool {
	return (*ContainerManager)(manager).AddConfig()
}

func (manager *PeerManager) Run() bool {
	log.Logger.Info("peer container start to create")
	return (*ContainerManager)(manager).Run()
}

func (manager *PeerManager) Restart() (ok bool) {
	return (*ContainerManager)(manager).Restart()
}

func (manager *PeerManager) StopAndRemove() (ok bool) {
	log.Logger.Info("peer container start to stop and remove")
	return (*ContainerManager)(manager).StopAndRemove()
}

func (manager *PeerManager) Stop() (ok bool) {
	return (*ContainerManager)(manager).Stop()
}

func (manager *PeerManager) Remove() (ok bool) {
	return (*ContainerManager)(manager).Remove()
}

func (manager *PeerManager) Exec(cmd []string) (content string, ok bool) {
	return (*ContainerManager)(manager).Exec(cmd)
}

func (manager *PeerManager) GetStatus() ContainerManager {
	return (ContainerManager)(*manager)
}
