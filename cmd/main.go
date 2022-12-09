package main

import (
	design "design/internal"
	"design/internal/conf"
	"design/internal/http"
	"flag"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const (
	ver     = "0.0.1"
	appName = "design-main"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("%s [version: %s env: %+v] start", appName, ver, conf.Conf.Env.Env)
	setLog(conf.Conf.Log, appName)
	// service code for main logic
	srv := design.New(conf.Conf)
	httpSrv := http.New(conf.Conf.HTTPServer, srv)
	_, err := strconv.Atoi(strings.Split(conf.Conf.HTTPServer.Addr, ":")[1])
	if err != nil {
		panic(err)
	}

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			httpSrv.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}

func setLog(logConfig *conf.LogConfig, appName string) {
	// set log output type as json
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	if logConfig.FileName == "" && appName != "" {
		logConfig.FileName = appName
	}
	path := logConfig.Path + logConfig.FileName + logConfig.Suffix + ".log"
	writer, _ := rotatelogs.New(
		path,
		//rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationCount(logConfig.RotationCount),
		//rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		rotatelogs.WithRotationSize(int64(logConfig.RotationSizeInMB*1024*1024)),
	)
	log.SetOutput(writer)

	level := log.WarnLevel
	if logConfig != nil && logConfig.Level != "" {
		switch logConfig.Level {
		case "Panic":
			{
				level = log.PanicLevel
			}
		case "Fatal":
			{
				level = log.FatalLevel
			}
		case "Error":
			{
				level = log.ErrorLevel
			}
		case "Warn":
			{
				level = log.WarnLevel
			}
		case "Info":
			{
				level = log.InfoLevel
			}
		case "Debug":
			{
				level = log.DebugLevel
			}
		}
	}

	// set default log level
	log.SetLevel(level)
}
