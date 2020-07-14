package docker

import (
	"context"
	"devplat/src/config"
	"devplat/src/log"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"net"
	"strings"
	"time"
)

const (
	networkName = "chaincode-docker-devmode-default"

	peerName      = "peer"
	ordererName   = "orderer"
	couchdbName   = "couchdb"
	cliName       = "cli"
	chaincodeName = "chaincode"
)

type DockerManager struct {
	cli                *client.Client
	ctx                context.Context
	CouchDBManager     *CouchDBManager
	OrdererManager     *OrdererManager
	PeerManager        *PeerManager
	CliManager         *CliManager
	ChaincodeManager   *ChaincodeManager
	ip                 string
	network            map[string]*network.EndpointSettings
	timeout            time.Duration
	dockerRunIntervals time.Duration
}

var dm *DockerManager

func InitDockerManager() *DockerManager {
	initCouchDBManager()
	initOrdererManager()
	initPeerManager()
	initCliManager()
	initChaincodeManager()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	config.IP, _ = getIp()
	var timeout = time.Second * 10
	dm = &DockerManager{
		cli:                cli,
		ctx:                context.Background(),
		ip:                 config.IP,
		CouchDBManager:     couchDBManager,
		OrdererManager:     ordererManager,
		PeerManager:        peerManager,
		CliManager:         cliManager,
		ChaincodeManager:   chaincodeManager,
		network:            make(map[string]*network.EndpointSettings),
		timeout:            timeout,
		dockerRunIntervals: 2 * time.Second,
	}

	dm.getImageID()
	return dm
}

func (dm *DockerManager) ImageInsure(imageIDMap map[string]string) {
	if _, exist := imageIDMap[dm.CouchDBManager.Image]; !exist {
		dm.PullImage(dm.CouchDBManager.Image)
	}
	if _, exist := imageIDMap[dm.OrdererManager.Image]; !exist {
		dm.PullImage(dm.OrdererManager.Image)
	}
	if _, exist := imageIDMap[dm.PeerManager.Image]; !exist {
		dm.PullImage(dm.PeerManager.Image)
	}
	if _, exist := imageIDMap[dm.CliManager.Image]; !exist {
		dm.PullImage(dm.CliManager.Image)
	}
	if _, exist := imageIDMap[dm.ChaincodeManager.Image]; !exist {
		dm.PullImage(dm.ChaincodeManager.Image)
	}
}

func (dm *DockerManager) getImageID() {
	images, err := dm.ListImages()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	imageIDMap := make(map[string]string)
	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			imageIDMap[repoTag] = strings.Split(image.ID, ":")[1]
		}
	}
	dm.ImageInsure(imageIDMap)
	dm.CouchDBManager.ID = imageIDMap[dm.CouchDBManager.Image]
	dm.OrdererManager.ID = imageIDMap[dm.OrdererManager.Image]
	dm.PeerManager.ID = imageIDMap[dm.PeerManager.Image]
	dm.CliManager.ID = imageIDMap[dm.CliManager.Image]
	dm.ChaincodeManager.ID = imageIDMap[dm.ChaincodeManager.Image]
}

func (dm *DockerManager) SetupDockerManager() {
	log.Logger.Info("start to setup fabric dev environment")
	initNetwork()

	if ok := dm.CouchDBManager.Run(); !ok {
		log.Logger.Warn("couchdb container create failed")
	} else {
		log.Logger.Info("couchdb container create success")
	}
	time.Sleep(dm.dockerRunIntervals)

	if ok := dm.OrdererManager.Run(); !ok {
		log.Logger.Warn("orderer container create failed")
	} else {
		log.Logger.Info("orderer container create success")
	}
	time.Sleep(dm.dockerRunIntervals)

	if ok := dm.PeerManager.Run(); !ok {
		log.Logger.Warn("peer container create failed")
	} else {
		log.Logger.Info("peer container create success")
	}
	time.Sleep(dm.dockerRunIntervals)

	if ok := dm.CliManager.Run(); !ok {
		log.Logger.Warn("cli container create failed")
	} else {
		log.Logger.Info("cli container create success")
	}
	time.Sleep(dm.dockerRunIntervals)

	if ok := dm.ChaincodeManager.Run(); !ok {
		log.Logger.Warn("chaincode container create failed")
	} else {
		log.Logger.Info("chaincode container create success")
	}
}

func (dm *DockerManager) CleanDockerManager() {
	log.Logger.Info("start to clean fabric dev environment")
	if ok := dm.ChaincodeManager.StopAndRemove(); !ok {
		log.Logger.Warn("chaincode container stop and remove failed")
	} else {
		log.Logger.Info("chaincode container stop and remove success")
	}

	if ok := dm.CliManager.StopAndRemove(); !ok {
		log.Logger.Warn("cli container stop and remove failed")
	} else {
		log.Logger.Info("cli container stop and remove success")
	}

	if ok := dm.PeerManager.StopAndRemove(); !ok {
		log.Logger.Warn("peer container stop and remove failed")
	} else {
		log.Logger.Info("peer container stop and remove success")
	}

	if ok := dm.OrdererManager.StopAndRemove(); !ok {
		log.Logger.Warn("orderer container stop and remove failed")
	} else {
		log.Logger.Info("orderer container stop and remove success")
	}

	if ok := dm.CouchDBManager.StopAndRemove(); !ok {
		log.Logger.Warn("couchdb container stop and remove failed")
	} else {
		log.Logger.Info("couchdb container stop and remove success")
	}
	dm.deleteNetwork()
}

func (dm *DockerManager) CheckHealthy() {
	containers, err := dm.ListContainer()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var containerMap = make(map[string]struct{})
	for _, c := range containers {
		containerMap[c.ID] = struct{}{}
	}
	_, dm.CouchDBManager.Healthy = containerMap[dm.CouchDBManager.ContainerID]
	_, dm.PeerManager.Healthy = containerMap[dm.PeerManager.ContainerID]
	_, dm.OrdererManager.Healthy = containerMap[dm.OrdererManager.ContainerID]
	_, dm.CliManager.Healthy = containerMap[dm.CliManager.ContainerID]
	_, dm.ChaincodeManager.Healthy = containerMap[dm.ChaincodeManager.ContainerID]
}

func (dm *DockerManager) deleteNetwork() (ok bool) {
	networkID := dm.network[networkName].NetworkID
	if err := dm.DeleteNetwork(networkID); err != nil {
		log.Logger.Error(err.Error())
		return false
	}
	delete(dm.network, networkName)
	return true
}

func initNetwork() {
	log.Logger.Info("start to init network")
	networks, err := dm.ListNetwork()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, nw := range networks {
		if nw.Name == networkName {
			dm.network[networkName] = &network.EndpointSettings{
				NetworkID: nw.ID,
			}
			return
		}
	}
	dm.CreateNetwork(networkName)
	log.Logger.Info(fmt.Sprintf("network %v created success", networkName))
}

func getIp() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return
}

func (dm *DockerManager) ClientVersion() string {
	return dm.cli.ClientVersion()
}

func (dm *DockerManager) CreateNetwork(name string) {
	resp, err := dm.cli.NetworkCreate(dm.ctx, name, types.NetworkCreate{})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	dm.network[name] = &network.EndpointSettings{
		NetworkID: resp.ID,
	}
}

func (dm *DockerManager) DeleteNetwork(networkID string) error {
	return dm.cli.NetworkRemove(dm.ctx, networkID)
}

func (dm *DockerManager) ListNetwork() (networks []types.NetworkResource, err error) {
	return dm.cli.NetworkList(dm.ctx, types.NetworkListOptions{})
}

func (dm *DockerManager) PullImage(imageName string) {
	dm.cli.ImagePull(dm.ctx, imageName, types.ImagePullOptions{})
}

func (dm *DockerManager) ListImages() ([]types.ImageSummary, error) {
	return dm.cli.ImageList(dm.ctx, types.ImageListOptions{})
}

func (dm *DockerManager) ListContainer() (containers []types.Container, err error) {
	return dm.cli.ContainerList(dm.ctx, types.ContainerListOptions{})
}

func (dm *DockerManager) RunContainer(config *container.Config, hostConfig *container.HostConfig, networkConfig *network.NetworkingConfig, containerName string) (containerID string, ok bool) {
	c, err := dm.cli.ContainerCreate(dm.ctx, config, hostConfig, networkConfig, containerName)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err := dm.cli.ContainerStart(dm.ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return c.ID, true
}

func (dm *DockerManager) StopContainer(containerID string) (ok bool) {
	err := dm.cli.ContainerStop(dm.ctx, containerID, &dm.timeout)
	if err != nil {
		log.Logger.Error(err.Error())
		return false
	}
	return true
}

func (dm *DockerManager) RemoveContainer(containerID string) (ok bool) {
	err := dm.cli.ContainerRemove(dm.ctx, containerID, types.ContainerRemoveOptions{})
	if err != nil {
		log.Logger.Error(err.Error())
		return false
	}
	return true
}

func (dm *DockerManager) Exec(cmd []string, containerID string) (content string, ok bool) {
	var cfg = types.ExecConfig{
		Privileged:   true,
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	execID, err := dm.cli.ContainerExecCreate(dm.ctx, containerID, cfg)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	res, err := dm.cli.ContainerExecAttach(dm.ctx, execID.ID, types.ExecConfig{})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	err = dm.cli.ContainerExecStart(dm.ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var slice = make([]string, 0)
	for {
		line, err := res.Reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		slice = append(slice, string(line))
	}
	return strings.Join(slice, ""), true
}
