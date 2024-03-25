package api

import (
	"fmt"
	"net"
	"video/pkg/config"
	"video/pkg/pb/video"

	"google.golang.org/grpc"
)

type Server struct{
	gs *grpc.Server
	Lis net.Listener
	Port string
}

func NewGrpcServer(c *config.Config,servcie video.VideoServiceServer) (*Server,error) {
	lis,err:=net.Listen("tcp",c.Port)
	if err != nil{
		return nil,fmt.Errorf("failed to listen on port")
	}

	grpcServer:=grpc.NewServer()
	video.RegisterVideoServiceServer(grpcServer,servcie)

	return &Server{
		gs: grpcServer,
		Lis: lis,
		Port: c.Port,
	},nil
}

func (s *Server) Start() error{
	fmt.Println("Service start on",s.Port)
	return s.gs.Serve(s.Lis)
}