package adapters

import "github.com/vishnusunil243/UserService/entities"

type AdapterInterface interface {
	UserSignup(req entities.User) (entities.User, error)
	GetUser(id int) (entities.User, error)
	GetAdmin(id int) (entities.Admin, error)
	UserLogin(email string) (entities.User, error)
	AdminLogin(email string) (entities.Admin, error)
	SuperAdminLogin(email string) (entities.SuperAdmin, error)
	AddAdmin(req entities.Admin) (entities.Admin, error)
	GetAllAdmins() ([]entities.Admin, error)
	GetAllUsers() ([]entities.User, error)
}
