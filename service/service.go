package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
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
func (admin *UserService) GetAllUsers(em *empty.Empty, srv pb.UserService_GetAllUsersServer) error {
	span := Tracer.StartSpan("get all users")
	defer span.Finish()
	users, err := admin.Adapter.GetAllUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		if err = srv.Send(&pb.UserSignupResponse{
			Id:    uint32(user.Id),
			Name:  user.Name,
			Email: user.Email,
		}); err != nil {
			return err
		}
	}
	return nil
}
func (sup *UserService) GetAllAdmins(em *empty.Empty, srv pb.UserService_GetAllAdminsServer) error {
	span := Tracer.StartSpan("get all admins grpc")
	defer span.Finish()
	admins, err := sup.Adapter.GetAllAdmins()
	if err != nil {
		return err
	}
	for _, admin := range admins {
		fmt.Println(admin)
		if err = srv.Send(&pb.UserSignupResponse{
			Id:    uint32(admin.Id),
			Name:  admin.Name,
			Email: admin.Email,
		}); err != nil {
			return err
		}
	}
	return nil
}
func (sup *UserService) AddAdmin(ctx context.Context, req *pb.UserSignupRequest) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("add admin grpc")
	defer span.Finish()
	reqEntity := entities.Admin{
		Name:     req.Name,
		Email:    req.Email,
		Password: helper.Hash(req.Password),
	}
	admin, err := sup.Adapter.AddAdmin(reqEntity)
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	res := &pb.UserSignupResponse{
		Id:    uint32(admin.Id),
		Name:  admin.Name,
		Email: admin.Email,
	}
	return res, nil
}
func (admin *UserService) GetUser(ctx context.Context, req *pb.GetUserById) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("get user grpc")
	defer span.Finish()
	user, err := admin.Adapter.GetUser(int(req.Id))
	if err != nil {
		return &pb.UserSignupResponse{}, err
	}
	res := &pb.UserSignupResponse{
		Id:    uint32(user.Id),
		Name:  user.Name,
		Email: user.Email,
	}
	return res, nil
}
func (sup *UserService) GetAdmin(ctx context.Context, req *pb.GetUserById) (*pb.UserSignupResponse, error) {
	span := Tracer.StartSpan("get admin grpc")
	defer span.Finish()
	admin, err := sup.Adapter.GetAdmin(int(req.Id))
	if err != nil {
		return &pb.UserSignupResponse{}, nil
	}
	res := &pb.UserSignupResponse{
		Id:    uint32(admin.Id),
		Name:  admin.Name,
		Email: admin.Email,
	}
	return res, nil
}
func (address *UserService) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.GetUserById, error) {
	addr, err := address.Adapter.GetAddress(int(req.UserId))
	if addr.Id != 0 {
		return &pb.GetUserById{}, fmt.Errorf("user already has address please delete the address to add new one")
	}
	if err != nil {
		return &pb.GetUserById{}, err
	}
	entity := entities.Address{
		UserId:   uint(req.UserId),
		City:     req.City,
		State:    req.State,
		District: req.District,
		Road:     req.Road,
	}
	err = address.Adapter.AddAddress(entity)
	if err != nil {
		return &pb.GetUserById{}, err
	}
	return &pb.GetUserById{Id: req.UserId}, nil
}
func (address *UserService) RemoveAddress(ctx context.Context, req *pb.GetUserById) (*pb.GetUserById, error) {
	addr, err := address.Adapter.GetAddress(int(req.Id))
	if err != nil {
		return &pb.GetUserById{}, err
	}
	if addr.Id == 0 {
		return &pb.GetUserById{}, fmt.Errorf("address not found")
	}
	err = address.Adapter.RemoveAddress(int(req.Id))
	if err != nil {
		return &pb.GetUserById{}, err
	}
	return &pb.GetUserById{Id: req.Id}, nil
}
func (address *UserService) GetAddress(ctx context.Context, req *pb.GetUserById) (*pb.AddAddressRequest, error) {
	addr, err := address.Adapter.GetAddress(int(req.Id))
	if err != nil {
		return &pb.AddAddressRequest{}, err
	}
	res := &pb.AddAddressRequest{
		UserId:   uint32(addr.UserId),
		City:     addr.City,
		State:    addr.State,
		Road:     addr.Road,
		District: addr.District,
	}
	fmt.Println(res)
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
