//go:build wireinject

package startup

import (
	grpc2 "gindemo/webook/comment/grpc"
	"gindemo/webook/comment/ioc"
	"gindemo/webook/comment/repository"
	"gindemo/webook/comment/repository/dao"
	"gindemo/webook/comment/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO, repository.NewCommentRepo,
	service.NewCommentSvc, grpc2.NewGrpcServer)

var thirdProvider = wire.NewSet(
	ioc.InitLogger, InitTestDB)

func InitGRPCServer() *grpc2.CommentServiceServer {
	wire.Build(thirdProvider, serviceProviderSet)
	return new(grpc2.CommentServiceServer)
}
