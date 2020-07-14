package docker

import (
	"devplat/src/log"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

type Port struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type ContainerManagerInterface interface {
	//给容器添加参数
	AddConfig() bool
	//运行容器
	Run() bool
	//重启容器
	Restart() bool
	//删除并停止容器
	StopAndRemove() bool
	//停止容器运行
	Stop() bool
	//删除容器
	Remove() bool
	//容器命令执行
	Exec(cmd []string) (content string, ok bool)
	//查看容器信息
	GetStatus() ContainerManager
}

type ContainerManager struct {
	ID            string            `json:"id"`
	Image         string            `json:"image"`
	ContainerID   string            `json:"container_id"`
	ContainerName string            `json:"container_name"`
	Env           []string          `json:"-"`
	Cmd           []string          `json:"-"`
	Volumes       []string          `json:"volumes"`
	WorkingDir    string            `json:"working_dir"`
	Ports         map[string][]Port `json:"ports"`
	Healthy       bool              `json:"healthy"`
	config        *container.Config
	hostConfig    *container.HostConfig
	networkConfig *network.NetworkingConfig
}

/*
	添加参数
*/
func (manager *ContainerManager) AddConfig() (ok bool) {
	var cfg = &container.Config{
		Image:        manager.Image,
		Env:          manager.Env,
		ExposedPorts: make(map[nat.Port]struct{}),
		WorkingDir:   manager.WorkingDir,
		Tty:          true,
		Cmd:          manager.Cmd,
	}
	var hostConfig = &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Binds:        manager.Volumes,
		PortBindings: make(map[nat.Port][]nat.PortBinding),
		Privileged:   true,
	}
	for extenral, exports := range manager.Ports {
		port, err := nat.NewPort("tcp", extenral)
		if err != nil {
			log.Logger.Error(err.Error())
			return false
		}
		cfg.ExposedPorts[port] = struct{}{}
		var portBindings = make([]nat.PortBinding, 0)
		for _, export := range exports {
			portBindings = append(portBindings, nat.PortBinding{
				HostIP:   export.IP,
				HostPort: export.Port,
			})
		}
		hostConfig.PortBindings[port] = portBindings
	}
	var networkConfig = &network.NetworkingConfig{
		EndpointsConfig: dm.network,
	}

	manager.config = cfg
	manager.hostConfig = hostConfig
	manager.networkConfig = networkConfig
	return true
}

func (manager *ContainerManager) PullImage() bool {
	return true
}

func (manager *ContainerManager) Run() bool {
	if manager.AddConfig() {
		containerID, ok := dm.RunContainer(manager.config, manager.hostConfig, manager.networkConfig, manager.ContainerName)
		manager.ContainerID = containerID
		return ok
	}
	return false
}

func (manager *ContainerManager) Restart() (ok bool) {
	if manager.StopAndRemove() {
		return manager.Run()
	}
	return false
}

func (manager *ContainerManager) StopAndRemove() (ok bool) {
	if dm.StopContainer(manager.ContainerID) {
		return dm.RemoveContainer(manager.ContainerID)
	}
	return false
}

func (manager *ContainerManager) Stop() (ok bool) {
	return dm.StopContainer(manager.ContainerID)
}

func (manager *ContainerManager) Remove() (ok bool) {
	return dm.RemoveContainer(manager.ContainerID)
}

func (manager *ContainerManager) Exec(cmd []string) (content string, ok bool) {
	return dm.Exec(cmd, manager.ContainerID)
}

func (manager *ContainerManager) GetStatus() ContainerManager {
	return *manager
}
