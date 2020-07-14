package docker

import "devplat/src/log"

type CouchDBManager ContainerManager

var couchDBManager *CouchDBManager

func initCouchDBManager() {
	var containerName = couchdbName
	couchDBManager = &CouchDBManager{
		Env:           []string{},
		ContainerName: containerName,
		Image:         "hyperledger/fabric-couchdb:latest",
		Ports: map[string][]Port{
			"5984": {
				{
					IP:   "0.0.0.0",
					Port: "5984",
				},
			},
		},
	}
}

func (manager *CouchDBManager) AddConfig() bool {
	return (*ContainerManager)(manager).AddConfig()
}

func (manager *CouchDBManager) Run() bool {
	log.Logger.Info("couchdb container start to create")
	return (*ContainerManager)(manager).Run()
}

func (manager *CouchDBManager) Restart() (ok bool) {
	return (*ContainerManager)(manager).Restart()
}

func (manager *CouchDBManager) StopAndRemove() (ok bool) {
	log.Logger.Info("couchdb container start to stop and remove")
	return (*ContainerManager)(manager).StopAndRemove()
}

func (manager *CouchDBManager) Stop() (ok bool) {
	return (*ContainerManager)(manager).Stop()
}

func (manager *CouchDBManager) Remove() (ok bool) {
	return (*ContainerManager)(manager).Remove()
}

func (manager *CouchDBManager) Exec(cmd []string) (content string, ok bool) {
	return (*ContainerManager)(manager).Exec(cmd)
}

func (manager *CouchDBManager) GetStatus() ContainerManager {
	return (ContainerManager)(*manager)
}
