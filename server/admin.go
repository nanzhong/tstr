package server

import "github.com/nanzhong/tstr/api/admin/v1"

type AdminServer struct {
	admin.UnimplementedAdminServiceServer
}

func NewAdminServer() admin.AdminServiceServer {
	return AdminServer{}
}
