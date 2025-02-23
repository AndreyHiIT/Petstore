package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	AppName = "Pet_Strore"

	serverPort         = "SERVER_PORT"
	levelLogger        = "LEVEL_LOGGER"
	envShutdownTimeout = "SHUTDOWN_TIMEOUT"
	envAccessTTL       = "ACCESS_TTL"
	envRefreshTTL      = "REFRESH_TTL"
	envVerifyLinkTTL   = "VERIFY_LINK_TTL"

	parseShutdownTimeoutError    = "config: parse server shutdown timeout error"
	parseRpcShutdownTimeoutError = "config: parse rpc server shutdown timeout error"
	parseTokenTTlError           = "config: parse token ttl error"
)

type AppConf struct {
	AppName         string `yaml:"app_name"`
	Environment     string `yaml:"environment"`
	Domain          string `yaml:"domain"`
	APIUrl          string `yaml:"api_url"`
	Server          Server `yaml:"server"`
	Token           Token  `yaml:"token"`
	Logger          Logger `yaml:"logger"`
	DB              DB     `yaml:"db"`
}

type Token struct {
	AccessTTL     time.Duration `yaml:"access_ttl"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl"`
	AccessSecret  string        `yaml:"access_secret"`
	RefreshSecret string        `yaml:"refresh_secret"`
}

type Server struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DB struct {
	Net      string `yaml:"net"`
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	User     string `json:"-" yaml:"user"`
	Password string `json:"-" yaml:"password"`
	Host     string `yaml:"host"`
	MaxConn  int    `yaml:"max_conn"`
	Port     string `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
}

type Logger struct {
	Level string `yaml:"level"`
}

func NewAppConf() AppConf {
	port := os.Getenv(serverPort)

	return AppConf{
		AppName:         os.Getenv(AppName),
		Server: Server{
			Port: port,
		},
		DB: DB{
			Net:      os.Getenv("DB_NET"),
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
	}
}

func (a *AppConf) Init(logger *zap.Logger) {
	shutDownTimeOut, err := strconv.Atoi(os.Getenv(envShutdownTimeout))
	if err != nil {
		logger.Fatal(parseShutdownTimeoutError)
	}
	shutDownTimeout := time.Duration(shutDownTimeOut) * time.Second

	dbTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse db timeout err", zap.Error(err))
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("MAX_CONN"))
	if err != nil {
		logger.Fatal("config: parse db max connection err", zap.Error(err))
	}
	a.DB.Timeout = dbTimeout
	a.DB.MaxConn = dbMaxConn

	var accessTTL int
	accessTTL, err = strconv.Atoi(os.Getenv(envAccessTTL))
	if err != nil {
		logger.Fatal(parseTokenTTlError)
	}
	a.Token.AccessTTL = time.Duration(accessTTL) * time.Minute
	var refreshTTL int
	refreshTTL, err = strconv.Atoi(os.Getenv(envRefreshTTL))
	if err != nil {
		logger.Fatal(parseTokenTTlError)
	}

	a.Token.AccessSecret = os.Getenv("ACCESS_SECRET")
	a.Token.RefreshSecret = os.Getenv("REFRESH_SECRET")
	a.Domain = os.Getenv("DOMAIN")
	a.APIUrl = os.Getenv("API_URL")

	a.Token.RefreshTTL = time.Duration(refreshTTL) * time.Hour * 24

	a.Server.ShutdownTimeout = shutDownTimeout
}
