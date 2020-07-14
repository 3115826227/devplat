package service

import (
	"devplat/src/service/docker"
	"sort"
)

var controller *DevPlatController

type DevPlatController struct {
	dockerManager *docker.DockerManager
	containers    map[string]docker.ContainerManagerInterface
	status        bool
}

func InitDevPlatController() {
	controller = &DevPlatController{
		dockerManager: docker.InitDockerManager(),
		containers:    make(map[string]docker.ContainerManagerInterface),
	}
}

func GetDevPlatController() *DevPlatController {
	return controller
}

func GetDockerManager() *docker.DockerManager {
	return controller.dockerManager
}

/*
	开启环境
*/
func (controller *DevPlatController) Setup() {
	if controller.status {
		return
	}
	controller.dockerManager.SetupDockerManager()
	dockerManager := controller.dockerManager
	controller.containers[dockerManager.ChaincodeManager.GetStatus().ContainerName] = dockerManager.ChaincodeManager
	controller.containers[dockerManager.CliManager.GetStatus().ContainerName] = dockerManager.CliManager
	controller.containers[dockerManager.PeerManager.GetStatus().ContainerName] = dockerManager.PeerManager
	controller.containers[dockerManager.OrdererManager.GetStatus().ContainerName] = dockerManager.OrdererManager
	controller.containers[dockerManager.CouchDBManager.GetStatus().ContainerName] = dockerManager.CouchDBManager
	controller.status = true
}

/*
	检测容器运行情况
*/
func (controller *DevPlatController) checkHealthy() {
	controller.dockerManager.CheckHealthy()
}

/*
	查看容器
*/
func (controller *DevPlatController) GetContainers() (status bool, containers []docker.ContainerManager) {
	containers = make([]docker.ContainerManager, 0)
	status = controller.status
	if !status {
		return
	}
	controller.checkHealthy()
	for _, manager := range controller.containers {
		containers = append(containers, manager.GetStatus())
	}
	sort.Sort(docker.Managers(containers))
	return
}

//
///*
//	重启容器
//*/
//func (controller *DevPlatController) RestartContainer(containerID string) {
//	controller.containers[containerID].Restart()
//	newContainer := controller.containers[containerID]
//	newContainerID := newContainer.GetStatus().ContainerID
//	controller.containers[newContainerID] = newContainer
//	delete(controller.containers, containerID)
//}

/*
	清除环境
*/
func (controller *DevPlatController) Clean() {
	if controller.status {
		controller.dockerManager.CleanDockerManager()
		controller.containers = make(map[string]docker.ContainerManagerInterface)
		controller.status = false
	}
}
