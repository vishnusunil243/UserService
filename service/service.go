package service

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/vishnusunil243/UserService/adapters"
	"github.com/vishnusunil243/UserService/entities"
	"github.com/vishnusunil243/UserService/helper"
	"github.com/vishnusunil243/proto-files/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	Tracer opentracing.Tracer
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type UserService struct {
	Adapter adapters.AdapterInterface
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapter adapters.AdapterInterface) *UserService {
	return &UserService{
		Adapter: adapter,
	}
}
func (user *UserService) UserSignup(ctx context.Context, req *pb.UserSignupRequest) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("user signup grpc")
	defer span.Finish()
	reqEntity := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: helper.Hash(req.Password),
	}
	res, err := user.Adapter.UserSignup(reqEntity)
	if err != nil {
		return nil, err
	}
	return &pb.UserSignupResponse{
		Id:    uint32(res.Id),
		Name:  res.Name,
		Email: res.Email,
	}, nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Println("check called")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
