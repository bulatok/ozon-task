package v1

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	"github.com/bulatok/ozon-task/internal/ozon-task/usecase"
)

type UseCases struct {
	Links *usecase.Links
}

func Init(grpcSrv *grpc.Server, l *zap.Logger, conf *config.Config, uc UseCases) {
	{
		initLinksService(grpcSrv, l, conf, uc.Links)
	}
}
