package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (svc *PinichleService) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	var status = grpc_health_v1.HealthCheckResponse_SERVING

	log.WithField("status", status).Debugf("health check: %s", status)

	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}, nil
}

func (s *PinichleService) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	var status = grpc_health_v1.HealthCheckResponse_SERVING

	log.WithField("status", status).Debugf("health check: %s", status)

	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: status,
	})
}
