package adapters

import "github.com/vishnusunil243/UserService/entities"

type AdapterInterface interface {
	UserSignup(req entities.User) (entities.User, error)
}
