package server

import "github.com/nanzhong/tstr/api/control/v1"

type ControlServer struct {
	control.UnimplementedControlServiceServer
}

func NewControlServer() control.ControlServiceServer {
	return &ControlServer{}
}
