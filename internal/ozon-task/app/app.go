package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	grpc_v1 "github.com/bulatok/ozon-task/internal/ozon-task/delivery/grpc/v1"
	http_v1 "github.com/bulatok/ozon-task/internal/ozon-task/delivery/http/api/v1"
	"github.com/bulatok/ozon-task/internal/ozon-task/store"
	"github.com/bulatok/ozon-task/internal/ozon-task/store/cache"
	"github.com/bulatok/ozon-task/internal/ozon-task/store/postgres"
	"github.com/bulatok/ozon-task/internal/ozon-task/store/redis"
	"github.com/bulatok/ozon-task/internal/ozon-task/usecase"
	"github.com/bulatok/ozon-task/pkg/logger"
)

func Run(conf *config.Config) {
	// logger
	l, err := logger.ProvideZap(conf)
	if err != nil {
		log.Fatal("could not configure logger: ", err)
	}
	l.Info("starting the service")

	// main context
	_ = context.Background()

	// repo
	linksRepo, err := provideLinksRepo(conf)
	if err != nil {
		l.Fatal("could not initialize store", zap.String("error", err.Error()))
	}
	l.Info("connected to repository", zap.String("type", linksRepo.Name()))

	// use cases
	linksUseCase := usecase.ProvideLinks(linksRepo, l)

	// http api v1
	httpApp := fiber.New(fiber.Config{
		ErrorHandler: http_v1.ErrorHandler,
	})

	srv := http_v1.ProvideServer(conf, l, http_v1.UseCases{
		Links: linksUseCase,
	})
	srv.Init(httpApp)

	go func() {
		l.Info("starting http server", zap.String("address", conf.HTTP.Host+":"+conf.HTTP.Port))
		if err := httpApp.Listen(conf.HTTP.Host + ":" + conf.HTTP.Port); err != nil {
			l.Fatal("could not start http server", zap.String("error", err.Error()))
		}
	}()

	// grpc v1
	lst, err := net.Listen("tcp", conf.Grpc.Host+":"+conf.Grpc.Port)
	if err != nil {
		l.Fatal("could not initialize tcp listener", zap.String("error", err.Error()))
	}

	grpcSrv := grpc.NewServer()
	grpc_v1.Init(grpcSrv, l, conf, grpc_v1.UseCases{
		Links: linksUseCase,
	})

	go func() {
		l.Info("starting grpc server", zap.String("address", conf.Grpc.Host+":"+conf.Grpc.Port))
		if err := grpcSrv.Serve(lst); err != nil {
			l.Fatal("could not start grpc server", zap.String("error", err.Error()))
		}
	}()

	l.Info("finished initializing the service")
	// graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	sig := <-signalCh

	// flushing everything
	l.Info("graceful shutdown", zap.String("signal", sig.String()))
	if err := linksRepo.Close(); err != nil {
		l.Info("could not close links repo",
			zap.String("store_name", linksRepo.Name()),
			zap.String("error", err.Error()))
	}

	if err := httpApp.Shutdown(); err != nil {
		l.Info("could not shutdown http server", zap.String("error", err.Error()))
	}

	if err := lst.Close(); err != nil {
		l.Info("could not close grpc server listener", zap.String("error", err.Error()))
	}
}

func provideLinksRepo(conf *config.Config) (store.LinksRepo, error) {
	switch conf.Service.StoreType {
	case config.RedisType:
		return redis.Provide(conf.Store.Redis)
	case config.PostgresType:
		return postgres.Provide(conf.Store.Postgres)
	case config.CacheType:
		return cache.Provide()
	}
	return nil, config.ErrUnknownStoreType
}
