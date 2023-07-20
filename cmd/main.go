package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/yaml.v2"
	"os"
	"prometheus-webhook-wechat/controller"
	"prometheus-webhook-wechat/models"
	"strings"
	"time"
)

var NotifyTargets []models.Target
var LogConfig = &promlog.Config{}
var MaxContentLength int
var (
	ListenPort = kingpin.Flag(
		"listen-port",
		"The port to listen for API interface,default: :80",
	).Default(":80").String()
	ConfigFile = kingpin.Flag(
		"config.file",
		"The full path of config file,default: config.yml",
	).Default("config.yml").ExistingFile()
	TemplateFile = kingpin.Flag(
		"template.file",
		"The full path of template file,default: template.tmpl",
	).Default("template.tmpl").ExistingFile()
)

func init() {

	flag.AddFlags(kingpin.CommandLine, LogConfig)
	kingpin.Version(version.Print("prometheus-webhook-wechat"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(LogConfig)
	level.Info(logger).Log("msg", "Starting prometheus-webhook-wechat")

	configBytes, err := os.ReadFile(*ConfigFile)
	if err != nil {
		level.Error(logger).Log("msg", "Load config file error: ", err)
		os.Exit(10)
	}
	var config models.Config
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		level.Error(logger).Log("msg", "Unmarshal config file error: ", err)
	}
	for _, v := range config.Targets {
		NotifyTargets = append(NotifyTargets, v)
	}
	MaxContentLength = config.MaxContentLength
}

func handler(handler models.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := new(models.NewGinContext)
		ctx.Context = c
		ctx.NotifyTargets = NotifyTargets
		ctx.MaxContentLength = MaxContentLength
		ctx.TemplateFile = TemplateFile
		ctx.Logger = promlog.New(LogConfig)
		handler(ctx)
	}
}

func main() {

	Logger := promlog.New(LogConfig)
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/wechat/webhook/send", handler(controller.CallWechatController))

	endless.DefaultReadTimeOut = 10 * time.Second
	endless.DefaultWriteTimeOut = 30 * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20

	srv := endless.NewServer(*ListenPort, router)

	err := srv.ListenAndServe()
	if err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			level.Error(Logger).Log("msg", err)
			os.Exit(10)
		}
	}

}
