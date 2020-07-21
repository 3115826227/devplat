package config

import (
	"devplat/src/log"
	"gopkg.in/yaml.v2"
	"os"
)

type OrdererOrgs struct {
	Name          string        `yaml:"Name"`
	Domain        string        `yaml:"Domain"`
	EnableNodeOUs bool          `yaml:"EnableNodeOUs"`
	Specs         []OrdererSpec `yaml:"Specs"`
}

type OrdererSpec struct {
	Hostname string `yaml:"Hostname"`
}

type PeerOrgs struct {
	Name          string       `yaml:"Name"`
	Domain        string       `yaml:"Domain"`
	EnableNodeOUs bool         `yaml:"EnableNodeOUs"`
	Template      PeerTemplate `yaml:"Template"`
	Users         PeerUsers    `yaml:"Users"`
}

type PeerTemplate struct {
	Count int `yaml:"Count"`
}

type PeerUsers struct {
	Count int `yaml:"Count"`
}

type CryptoConfig struct {
	OrdererOrgs []OrdererOrgs `yaml:"OrdererOrgs"`
	PeerOrgs    []PeerOrgs    `yaml:"PeerOrgs"`
}

var (
	cryptoCfg CryptoConfig
)

func init() {
	cryptoCfg = CryptoConfig{
		OrdererOrgs: []OrdererOrgs{
			{
				Name:          "Orderer",
				Domain:        "example.com",
				EnableNodeOUs: true,
				Specs: []OrdererSpec{
					{
						Hostname: "orderer",
					},
				},
			},
		},
		PeerOrgs: []PeerOrgs{
			{
				Name:          "Org1",
				Domain:        "org1.example.com",
				EnableNodeOUs: true,
				Template:      PeerTemplate{Count: 2},
				Users:         PeerUsers{Count: 1},
			},
			{
				Name:          "Org2",
				Domain:        "org2.example.com",
				EnableNodeOUs: true,
				Template:      PeerTemplate{Count: 2},
				Users:         PeerUsers{Count: 1},
			},
		},
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func write(filename string, data []byte) {
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		f, err = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	if _, err = f.Write(data); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

func GenerateConfig() {
	data, err := yaml.Marshal(cryptoCfg)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	write("crypto.yaml", data)
}
