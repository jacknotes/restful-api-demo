package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	app "github.com/jacknotes/restful-api-demo/apps"
	"github.com/jacknotes/restful-api-demo/apps/host/impl"
	"github.com/jacknotes/restful-api-demo/conf"
	"github.com/jacknotes/restful-api-demo/protocol"
	"github.com/spf13/cobra"
)

var (
	configType string
	confFile   string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Demo后端API服务",
	Long:  `Demo后端API服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 加载全局配置
		if err := loadGlobalConfig(configType); err != nil {
			return err
		}

		// 初始化日志
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		// 初始化服务层，IOC初始化
		if err := impl.Service.Init(); err != nil {
			return err
		}
		// 把服务实例注册给IOC层，如果有多个不同服务实例实现，可以在这里解耦
		app.Host = impl.Service

		// 启动服务后，需要处理的事件，没有做信号区分，以下4种信号我都认为需要致使中断服务的
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		// 创建服务对象
		srv := NewService(conf.C())

		// 等待程序退出，启动协程处理
		go srv.waitSign(ch)

		// 启动服务
		return srv.Start()

	},
}

// http service
// 需要整体配置，可能启动很多模块：http, grpc, crontab
type Service struct {
	conf *conf.Config
	http *protocol.HTTPService // protocol层接口
	grpc *protocol.GrpcService // grpc层接口
	log  logger.Logger
}

func NewService(conf *conf.Config) *Service {
	return &Service{
		conf: conf,
		http: protocol.NewHTTPService(),
		grpc: protocol.NewGrpcService(),
		log:  zap.L().Named("Service"),
	}
}

func (s *Service) Start() error {
	go s.grpc.Start()
	return s.http.Start()
}

// 当用户手动终止程序的时候，需要完成处理
func (s *Service) waitSign(sign chan os.Signal) {
	for sg := range sign { // 如果未接收到通道数据则会一直阻塞
		switch v := sg.(type) { //v的type是os.Signal
		default:
			s.grpc.Stop()

			//资源整理
			s.log.Infof("receive signal '%v', start graceful shutdown", v.String())
			if err := s.http.Stop(); err != nil {
				s.log.Errorf("graceful shutdown err: %s, force exit", err)
			}
			s.log.Infof("service stop complete")
			return
		}
	}
}

// config 为全局变量, 只需要load 即可设置全局配置
func loadGlobalConfig(configType string) error {
	// 配置加载
	switch configType {
	case "file":
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			return err
		}
	case "env":
		err := conf.LoadConfigFromEnv()
		if err != nil {
			return err
		}
	case "etcd":
		return errors.New("not implemented")
	default:
		return errors.New("unknown config type")
	}
	return nil
}

// log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)

	// 获取出日志配置对象
	lc := conf.C().Log

	// 解析配置的日志级别是否正确，并且设置日志的级别
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		// 解析失败，默认使用info级别
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		// 解析成功，使用用户配置的日志级别
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}

	zapConfig := zap.DefaultConfig()
	zapConfig.Level = level
	zapConfig.Files.RotateOnStartup = false
	switch lc.To {
	case conf.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "restful-api.log"
		zapConfig.Files.Path = lc.PathDir
	}

	// 日志格式，如果为json就配置为json，否则默认为文本
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}

	// 初始化全局Logger的配置
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}

	// 全局Logger初始化后就可以正常使用了
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	RootCmd.AddCommand(startCmd)
	RootCmd.PersistentFlags().StringVarP(&configType, "config_type", "t", "file", "the restful-api config type")
	RootCmd.PersistentFlags().StringVarP(&confFile, "config_file", "f", "etc/restful-api.toml", "the restful-api config file path")
}
