package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	"github.com/bulatok/ozon-task/internal/ozon-task/usecase"
	"github.com/bulatok/ozon-task/pkg/pb"
)

type linksService struct {
	pb.UnimplementedLinksServer
	core         *grpc.Server
	linksUseCase *usecase.Links
	l            *zap.Logger
	conf         *config.Config
}

func initLinksService(grpcSrv *grpc.Server, l *zap.Logger, conf *config.Config, linksUseCase *usecase.Links) {
	childL := l.With(zap.String("logger", "grpc server"))

	service := &linksService{
		core:         grpcSrv,
		linksUseCase: linksUseCase,
		l:            childL,
		conf:         conf,
	}
	pb.RegisterLinksServer(grpcSrv, service)
}

func (s *linksService) ShortLink(_ context.Context, req *pb.ShortLinkRequest) (*pb.ShortLinkResponse, error) {
	s.l.Debug("short link")
	short, err := s.linksUseCase.New(req.OriginalLink, s.conf.Service.PublicUrl)
	resp := &pb.ShortLinkResponse{ShortLink: short}
	return resp, err
}

func (s *linksService) GetOriginalLink(_ context.Context, req *pb.GetOriginalRequest) (*pb.GetOriginalResponse, error) {
	s.l.Debug("get original link")
	original, err := s.linksUseCase.Get(req.ShortLink)
	resp := &pb.GetOriginalResponse{OriginalLink: original}
	return resp, err
}
