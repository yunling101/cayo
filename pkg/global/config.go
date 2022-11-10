package global

import (
	"log"
	"sync"

	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
)

type HttpConfig struct {
	Enable bool   `yaml:"enable"`
	Listen string `yaml:"listen"`
}

type DataBaseConfig struct {
	DBHost string `yaml:"host"`
	DBPort int    `yaml:"port"`
	DBUser string `yaml:"username"`
	DBPass string `yaml:"password"`
	DBName string `yaml:"database"`
	Prefix string `yaml:"prefix"`
}

type EmailConfig struct {
	Enable   bool   `yaml:"enable"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type ChannelConfig struct {
	Email *EmailConfig `yaml:"email"`
}

type GlobalCfg struct {
	Debug        bool            `yaml:"debug"`
	Workspace    string          `yaml:"workspace"`
	DataBase     *DataBaseConfig `yaml:"database"`
	RpcServer    *HttpConfig     `yaml:"rpcserver"`
	WebServer    *HttpConfig     `yaml:"webserver"`
	Channel      *ChannelConfig  `yaml:"channel"`
	NodeMaxLimit int             `yaml:"node_max_limit"`
}

var (
	config     *GlobalCfg
	configLock = new(sync.RWMutex)
)

func Config() *GlobalCfg {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func LoadConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent.")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var v GlobalCfg
	err = yaml.Unmarshal([]byte(configContent), &v)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &v

	log.Println("read config file:", cfg, "successfully")
}
