package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "vicomova/internal/shared/infrastructure/data/mysql"
    interactionAppCmd "vicomova/internal/interaction/application/command"
    interactionAppQuery "vicomova/internal/interaction/application/query"
    interactionInfraMysql "vicomova/internal/interaction/infrastructure/persistence/mysql"
    interactionGrpc "vicomova/internal/interaction/interfaces/grpc"
    "vicomova/internal/shared/pkg/log"
    "vicomova/pkg/config"
    "vicomova/pkg/etcd"

    interactionservice "vicomova/third_party/kitex_gen/interaction/interactionservice"

    "github.com/cloudwego/kitex/pkg/klog"
)

var (
    configPath string
    port       int
)

func init() {
    flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
    flag.IntVar(&port, "port", 8890, "interaction service port")
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

    // Initialize repositories
    likeRepo := interactionInfraMysql.NewLikeRepository()
    commentRepo := interactionInfraMysql.NewCommentRepository()
    favoriteRepo := interactionInfraMysql.NewFavoriteRepository()

    // Initialize application services
    cmdSvc := interactionAppCmd.NewInteractionCommandService(likeRepo, commentRepo, favoriteRepo)
    querySvc := interactionAppQuery.NewInteractionQueryService(likeRepo, commentRepo, favoriteRepo)

    // Initialize gRPC handler
    handler := interactionGrpc.NewInteractionHandler(cmdSvc, querySvc)

    // Create Kitex server
    svr := interactionservice.NewServer(handler)

    // Register to etcd
    var registry *etcd.Registry
    if len(cfg.Etcd.Endpoints) > 0 {
        etcdClient, err := etcd.NewClient(&cfg.Etcd)
        if err != nil {
            log.Info.Printf("Failed to create etcd client: %v", err)
        } else {
            defer etcdClient.Close()
            registry = etcd.NewRegistry(etcdClient, "interaction", fmt.Sprintf("127.0.0.1:%d", port))
            if err := registry.Register(); err != nil {
                log.Info.Printf("Failed to register to etcd: %v", err)
            }
        }
    }

    // Graceful shutdown
    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig
        klog.Info("Shutting down interaction service...")
        if registry != nil {
            registry.Unregister()
        }
        svr.Stop()
    }()

    log.Info.Printf("Interaction service starting on 127.0.0.1:%d", port)
    if err := svr.Run(); err != nil {
        log.Error.Fatalf("Server error: %v", err)
    }
}
