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
func (user *UserService) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("user login grpc")
	defer span.Finish()

	client, err := user.Adapter.UserLogin(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	if client.Name == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CheckPassword(client.Password, req.Password) {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")

	}
	res := &pb.UserSignupResponse{
		Id:    uint32(client.Id),
		Name:  client.Name,
		Email: client.Email,
	}

	return res, nil

}
func (admin *UserService) AdminLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("admin login grpc")
	defer span.Finish()
	adminData, err := admin.Adapter.AdminLogin(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	if adminData.Name == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CheckPassword(adminData.Password, req.Password) {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	res := &pb.UserSignupResponse{
		Id:    uint32(adminData.Id),
		Name:  adminData.Name,
		Email: adminData.Email,
	}
	return res, nil
}
func (sup *UserService) SuperAdminLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("super admin login grpc")
	defer span.Finish()
	supData, err := sup.Adapter.SuperAdminLogin(req.Email)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	if supData.Name == "" {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	if !helper.CheckPassword(supData.Password, req.Password) {
		return &pb.UserSignupResponse{}, fmt.Errorf("invalid credentials")
	}
	res := &pb.UserSignupResponse{
		Id:    uint32(supData.Id),
		Name:  supData.Name,
		Email: supData.Email,
	}
	return res, nil
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
