package config

import (
	"encoding/json"
	"go-shop/utils"
	"io/ioutil"
	"strings"

	"go.uber.org/zap"
)

var jsonData map[string]interface{}

/**
 * 服务端配置
 */
type serverConfig struct {
	AppName   string
	LogLevel  string
	Port      string
	SessionID string //后台设置的session id
	Mode      string
}

var ServerConfig serverConfig

func initServer() {
	utils.SetStructByJSON(&ServerConfig, jsonData["server"].(map[string]interface{}))
}

/**
 * mysql配置
 */
type dataBaseConfig struct {
	Drive    string
	Port     string
	User     string
	Pwd      string
	Host     string
	Database string
	Charset  string
	URL      string
	SQLLog   bool
	MaxOpenConns  int
}

var DataBaseConfig dataBaseConfig

func initDataBase() {
	utils.SetStructByJSON(&DataBaseConfig, jsonData["database"].(map[string]interface{}))
	url := "{user}:{password}@tcp({host}:{port})/{database}?charset={charset}&parseTime=True&loc=Local"
	url = strings.Replace(url, "{database}", DataBaseConfig.Database, -1)
	url = strings.Replace(url, "{user}", DataBaseConfig.User, -1)
	url = strings.Replace(url, "{password}", DataBaseConfig.Pwd, -1)
	url = strings.Replace(url, "{host}", DataBaseConfig.Host, -1)
	url = strings.Replace(url, "{port}", DataBaseConfig.Port, -1)
	url = strings.Replace(url, "{charset}", DataBaseConfig.Charset, -1)
	DataBaseConfig.URL = url
}

/**
 * Redis 配置
 */
type redisConfig struct {
	NetWork string
	Addr    string
	Port    string
	Pwd     string
	Prefix  string // key value set(name,"davie")  name -> xxx_name
}

var RedisConfig redisConfig

func initRedis() {
	utils.SetStructByJSON(&RedisConfig, jsonData["redis"].(map[string]interface{}))
}

func initJson() {
	bytes, err := ioutil.ReadFile("./config/conf.json")
	if err != nil {
		utils.Logger.Panic("ReadFile ./config/conf.json error", zap.Any("error", err))
	}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		utils.Logger.Panic("json.Unmarshal error", zap.Any("error", err))
	}
}

func init() {
	initJson()
	initServer()
	initDataBase()
	initRedis()
}

/*func InitConfig() *serverConfig {
	file, err := os.OpenFile("./config/conf.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	decoder := json.NewDecoder(file)
	conf := serverConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err.Error)
	}
	return &conf
}*/
