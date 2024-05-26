//go:build wireinject

package wire

import (
	"gindemo/wire/repository"
	"gindemo/wire/repository/dao"
	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, InitDB, dao.NewUserDAO)
	return &repository.UserRepository{}
}
