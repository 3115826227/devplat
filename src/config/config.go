package config

import "github.com/spf13/viper"

var (
	WorkPath     string
	GoPath       string
	IP           string
	Config       ConfigStruct
	TransientMap = make(map[string][]byte)
)

type ConfigStruct struct {
	ReleaseMode   bool   `env:"RELEASE_MODE" default:"false"`
	ChainCodePath string `env:"CHAIN_CODE_PATH" required:"true"`
	SdkCfgPath    string `env:"SDK_CFG_PATH" required:"true"`
	OrgName       string `env:"ORG_NAME" required:"true"`
	UserName      string `env:"USER_NAME" required:"true"`

	PeerImage    string            `env:"PEER_IMAGE" required:"true"`
	OrdererImage string            `env:"ORDERER_IMAGE" required:"true"`
	CcenvImage   string            `env:"CCENV_IMAGE" required:"true"`
	CouchdbName  string            `env:"COUCHDB_IMAGE" required:"true"`
	ChannelName  string            `env:"CHANNEL_NAME"`
	Peers        []string          `env:"PEERS" required:"true"`
	Orderer      []string          `env:"ORDERER" required:"true"`
	TransientMap map[string]string `env:"TRANSIENT_MAP" default:""`
}

func setConfigYaml() {
	//设置配置文件的名字
	viper.SetConfigName("config")
	//设置配置文件读取路径
	viper.AddConfigPath("./etc") //idea跑的时候直接读取项目etc/目录
	//viper.AddConfigPath("/etc")  //部署到docker容器中挂载到/etc目录下
	//设置配置文件类型
	viper.SetConfigType("yaml")
}

func init() {
	setConfigYaml()
	//读取配置文件内容
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var c ConfigStruct
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}
	Config = c
	for key, value := range Config.TransientMap {
		TransientMap[key] = []byte(value)
	}
}
