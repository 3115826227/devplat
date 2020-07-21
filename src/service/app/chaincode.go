package app

import (
	"devplat/src/config"
	"devplat/src/log"
	"io/ioutil"
	"os"
)

type ChaincodeInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ChaincodeProvider struct {
	status                 bool
	uninstallChaincodes    []string
	installedChaincodes    []ChaincodeInfo
	instantiatedChaincodes []string
}

var (
	chaincodeProvider *ChaincodeProvider
)

func InitChaincodeProvider() {
	chaincodeProvider = &ChaincodeProvider{
		uninstallChaincodes:    make([]string, 0),
		installedChaincodes:    make([]ChaincodeInfo, 0),
		instantiatedChaincodes: make([]string, 0),
	}
}

func findAllChaincodeDir() (subDirs []string, err error) {
	var dir = config.WorkPath + "/deploy/chaincode/"
	var fileInfo []os.FileInfo
	fileInfo, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range fileInfo {
		if f.IsDir() {
			if f.Name()[0] != '.' {
				subDirs = append(subDirs, f.Name())
			}
		}
	}
	return
}

func GetChaincodeProvider() *ChaincodeProvider {
	return chaincodeProvider
}

func (provider *ChaincodeProvider) Start() {
	provider.status = true
}

func (provider *ChaincodeProvider) GetUninstallChaincode() []string {
	if !provider.status {
		return nil
	}
	var installedChaincodeMap = make(map[string]struct{})
	var instantiatedChaincodeMap = make(map[string]struct{})
	for _, cc := range provider.installedChaincodes {
		installedChaincodeMap[cc.Name] = struct{}{}
	}
	for _, ccName := range provider.instantiatedChaincodes {
		instantiatedChaincodeMap[ccName] = struct{}{}
	}
	allChaincodes, err := findAllChaincodeDir()
	if err != nil {
		log.Logger.Error(err.Error())
		return provider.uninstallChaincodes
	}
	provider.uninstallChaincodes = make([]string, 0)
	for _, cc := range allChaincodes {
		if _, exist := installedChaincodeMap[cc]; exist {
			continue
		}
		if _, exist := instantiatedChaincodeMap[cc]; exist {
			continue
		}
		provider.uninstallChaincodes = append(provider.uninstallChaincodes, cc)
	}
	return provider.uninstallChaincodes
}

func (provider *ChaincodeProvider) ChaincodeInstallFeedback(chaincode ChaincodeInfo) {
	var uninstallChaincodes = make([]string, 0)
	for _, ccName := range provider.uninstallChaincodes {
		if ccName == chaincode.Name {
			provider.installedChaincodes = append(provider.installedChaincodes, chaincode)
			continue
		}
		uninstallChaincodes = append(uninstallChaincodes, ccName)
	}
	provider.uninstallChaincodes = uninstallChaincodes
}

func (provider *ChaincodeProvider) JudgeChaincodeInstalled(chaincode ChaincodeInfo) bool {
	for _, cc := range provider.installedChaincodes {
		if cc.Name == chaincode.Name && cc.Version == chaincode.Version {
			return true
		}
	}
	return false
}

func (provider *ChaincodeProvider) GetInstantiatedChaincode() []string {
	if !provider.status {
		return nil
	}
	return provider.instantiatedChaincodes
}

func (provider *ChaincodeProvider) ChaincodeInstantiateFeedback(chaincode ChaincodeInfo) {
	var installedChaincodes = make([]ChaincodeInfo, 0)
	for _, cc := range provider.installedChaincodes {
		if cc.Name == chaincode.Name && cc.Version == chaincode.Version {
			provider.instantiatedChaincodes = append(provider.instantiatedChaincodes, cc.Name)
			continue
		}
		installedChaincodes = append(installedChaincodes, cc)
	}
	provider.installedChaincodes = installedChaincodes
}

func (provider *ChaincodeProvider) JudgeChaincodeInstantiate(chaincodeName string) bool {
	for _, ccName := range provider.instantiatedChaincodes {
		if ccName == chaincodeName {
			return true
		}
	}
	return false
}

func (provider *ChaincodeProvider) GetInstalledChaincodes() []ChaincodeInfo {
	if !provider.status {
		return nil
	}
	return provider.installedChaincodes
}

func (provider *ChaincodeProvider) Clean() {
	provider.status = false
	provider.uninstallChaincodes = make([]string, 0)
	provider.installedChaincodes = make([]ChaincodeInfo, 0)
	provider.instantiatedChaincodes = make([]string, 0)
}
