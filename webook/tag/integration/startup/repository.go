package startup

import (
	"gindemo/webook/pkg/logger"
	"gindemo/webook/tag/repository"
	"gindemo/webook/tag/repository/cache"
	"gindemo/webook/tag/repository/dao"
)

func InitRepository(d dao.TagDAO, c cache.TagCache, l logger.LoggerV1) repository.TagRepository {
	return repository.NewTagRepository(d, c, l)
}
