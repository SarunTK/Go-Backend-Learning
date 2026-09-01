package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const serviceName = "user.v1.UserService"

type CreateUserRequest struct {
	Name  string
	Email string
}

type GetUserRequest struct {
	ID int64
}

type UserResponse struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt string
}

type ListUsersRequest struct{}

type ListUsersResponse struct {
	Users []*UserResponse
}

type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*UserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*UserResponse, error)
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
}

type serviceServer struct{}

func (s *serviceServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	return &UserResponse{
		ID:        1,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: "2026-01-01T00:00:00Z",
	}, nil
}

func (s *serviceServer) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	return &UserResponse{ID: req.ID, Name: "test101", Email: "test101@test101.com", CreatedAt: "2026-01-01T00:00:00Z"}, nil
}

func (s *serviceServer) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{Users: []*UserResponse{{ID: 1, Name: "test101", Email: "test101@test101.com", CreatedAt: "2026-01-01T00:00:00Z"}}}, nil
}

type Server struct {
	server *grpc.Server
}

func NewServer() *Server {
	s := &Server{server: grpc.NewServer()}
	s.server.RegisterService(&grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*UserServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "CreateUser",
				Handler:    createUserHandler,
			},
			{
				MethodName: "GetUser",
				Handler:    getUserHandler,
			},
			{
				MethodName: "ListUsers",
				Handler:    listUsersHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/user.proto",
	}, &serviceServer{})
	return s
}

func createUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/" + serviceName + "/CreateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/" + serviceName + "/GetUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func listUsersHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/" + serviceName + "/ListUsers"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ListUsers(ctx, req.(*ListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func (s *Server) ListenAndServe(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}
	log.Println("gRPC listener started")
	return s.server.Serve(ln)
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
