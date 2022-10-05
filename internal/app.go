package internal

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/costa92/go-web/config"
	"github.com/costa92/go-web/internal/db"
	"github.com/costa92/go-web/internal/option"
	"github.com/costa92/go-web/internal/options"
	genericapiserver "github.com/costa92/go-web/internal/server"
	"github.com/costa92/go-web/pkg/app"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/shutdown"
	"github.com/costa92/go-web/pkg/shutdown/shutdownmanagers/posixsignal"
)

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("go-web Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription("The go-web"),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		// 初始化日志
		logger.Init(opts.Logger)
		defer logger.Flush()
		// 配置文件赋值
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		return Run(cfg)
	}
}

// 绑定配置
func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	// 实例化config
	genericConfig = genericapiserver.NewConfig()
	// 加载配置文件
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = cfg.InsecureServingOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = cfg.Jwt.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	return
}

// 创建api服务
func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// 参数bind
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 参数重新赋值
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
	}
	return server, nil
}

type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   option.GeneratableKeyCert
	mysqlOptions *option.MySQLOptions
}

func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         "",
		ServerCert:   cfg.SecureServing.ServerCert,
		mysqlOptions: cfg.MysqlConfig,
	}, nil
}

type completedExtraConfig struct {
	*ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	var grpcServer *grpc.Server
	if c.ServerCert.CertKey.CertFile != "" && c.ServerCert.CertKey.KeyFile != "" {
		// 运行GRPc
		creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
		if err != nil {
			logger.Infof("Failed to generate credentials %s", err.Error())
		} else {
			opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
			grpcServer = grpc.NewServer(opts...)
			// pb.RegisterCacheServer(grpcServer, cacheIns)
		}
	}

	_, _ = db.GetMySQLFactoryOr(c.mysqlOptions)
	return &grpcAPIServer{grpcServer, c.Addr}, nil
}
