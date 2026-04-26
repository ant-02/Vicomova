package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "vicomova/internal/shared/infrastructure/data/mysql"
    videoAppCmd "vicomova/internal/video/application/command"
    videoAppQuery "vicomova/internal/video/application/query"
    videoInfraMysql "vicomova/internal/video/infrastructure/persistence/mysql"
    videoInfraService "vicomova/internal/video/infrastructure/service"
    videoInfraStorage "vicomova/internal/video/infrastructure/storage"
    videoGrpc "vicomova/internal/video/interfaces/grpc"
    "vicomova/internal/shared/pkg/log"
    "vicomova/pkg/config"
    "vicomova/pkg/etcd"

    videoservice "vicomova/third_party/kitex_gen/video/videoservice"

    "github.com/cloudwego/kitex/pkg/klog"
)

var (
    configPath string
    port       int
)

func init() {
    flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
    flag.IntVar(&port, "port", 8889, "video service port")
}

func main() {
    flag.Parse()

    klog.SetLevel(klog.LevelInfo)

    cfg, err := config.Load(configPath)
    if err != nil {
        log.Error.Fatalf("Failed to load config: %v", err)
    }

    // Initialize database
    if err := mysql.Init(&cfg.Database); err != nil {
        log.Error.Fatalf("Failed to init mysql: %v", err)
    }

    // Initialize repository
    videoRepo := videoInfraMysql.NewVideoRepository()

    // Initialize storage
    var storage videoInfraStorage.VideoStorage
    switch cfg.Video.Storage.Type {
    case "local":
        storage = videoInfraStorage.NewLocalDiskStorage(&videoInfraStorage.LocalDiskConfig{
            BasePath: cfg.Video.Storage.Local.BasePath,
            BaseURL:  cfg.Video.Storage.Local.BaseURL,
        })
    default:
        storage = videoInfraStorage.NewLocalDiskStorage(&videoInfraStorage.LocalDiskConfig{
            BasePath: cfg.Video.Storage.Local.BasePath,
            BaseURL:  cfg.Video.Storage.Local.BaseURL,
        })
    }

    // Initialize services
    hotAlgo := videoInfraService.NewWilsonHotAlgorithm(videoRepo)

    // Initialize application services
    cmdSvc := videoAppCmd.NewVideoCommandService(videoRepo, storage)
    querySvc := videoAppQuery.NewVideoQueryService(videoRepo, hotAlgo)

    // Initialize gRPC handler
    handler := videoGrpc.NewVideoHandler(cmdSvc, querySvc)

    // Create Kitex server
    svr := videoservice.NewServer(handler)

    // Register to etcd
    var registry *etcd.Registry
    if len(cfg.Etcd.Endpoints) > 0 {
        etcdClient, err := etcd.NewClient(&cfg.Etcd)
        if err != nil {
            log.Warn.Printf("Failed to create etcd client: %v", err)
        } else {
            defer etcdClient.Close()
            registry = etcd.NewRegistry(etcdClient, "video", fmt.Sprintf("127.0.0.1:%d", port))
            if err := registry.Register(); err != nil {
                log.Warn.Printf("Failed to register to etcd: %v", err)
            }
        }
    }

    // Graceful shutdown
    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig
        klog.Info("Shutting down video service...")
        if registry != nil {
            registry.Unregister()
        }
        svr.Stop()
    }()

    log.Info.Printf("Video service starting on 127.0.0.1:%d", port)
    if err := svr.Run(); err != nil {
        log.Error.Fatalf("Server error: %v", err)
    }
}
