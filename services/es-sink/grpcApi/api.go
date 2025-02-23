package grpcApi

import (
	"context"
	"encoding/json"
	"github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/proto/src/golang"
	"github.com/opengovern/opensecurity/pkg/utils"
	"github.com/opengovern/opensecurity/services/es-sink/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCSinkServer struct {
	golang.EsSinkServiceServer
	logger        *zap.Logger
	esSinkService *service.EsSinkService
	grpcServer    *grpc.Server
	addr          string
}

func NewGRPCSinkServer(logger *zap.Logger, esSinkService *service.EsSinkService, addr string) (*GRPCSinkServer, error) {
	server := GRPCSinkServer{
		logger:        logger,
		esSinkService: esSinkService,
		addr:          addr,
	}

	server.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(128 * 1024 * 1024),
	)
	golang.RegisterEsSinkServiceServer(server.grpcServer, &server)

	return &server, nil
}

func (s *GRPCSinkServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Fatal("failed to listen on grpc port", zap.Error(err))
	}
	utils.EnsureRunGoroutine(func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			s.logger.Fatal("failed to serve grpc server", zap.Error(err))
		}
	})
	return nil
}

func (s *GRPCSinkServer) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *GRPCSinkServer) Ingest(ctx context.Context, req *golang.IngestRequest) (*golang.ResponseOK, error) {
	docs := make([]es.DocBase, 0, len(req.Docs))
	for _, doc := range req.Docs {
		var d es.DocBase
		err := json.Unmarshal(doc.Value, &d)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	if _, err := s.esSinkService.Ingest(ctx, docs); err != nil {
		s.logger.Error("failed to ingest data", zap.Error(err))
		return nil, err
	}
	return &golang.ResponseOK{}, nil
}
