package initializer

import (
	"github.com/vishnusunil243/UserService/adapters"
	"github.com/vishnusunil243/UserService/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.UserService {
	adapter := adapters.NewUserAdapter(db)
	service := service.NewUserService(adapter)
	return service
}
