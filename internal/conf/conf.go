package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var (
	confPath string
	env      string
	Conf     *Config
)

type Config struct {
	Env        *Env
	HTTPServer *HTTPServer
	Mongodb    *Mongodb
	Log        *LogConfig
}

type Env struct {
	Env string
}

// Mongodb .
type Mongodb struct {
	Host                 string
	Username             string
	Password             string
	DatabaseName         string
	MaxPoolSize          uint64
	TimeoutInMillisecond int64
}

type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type LogConfig struct {
	Level            string
	Path             string
	FileName         string
	Suffix           string
	RotationCount    uint
	RotationSizeInMB uint
	EnableHook       bool
}

func init() {
	flag.StringVar(&confPath, "conf", "design.toml", "default config path")
	flag.StringVar(&env, "env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
