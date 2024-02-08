package adapters

import (
	"github.com/vishnusunil243/UserService/entities"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func NewUserAdapter(db *gorm.DB) *UserAdapter {
	return &UserAdapter{DB: db}
}
func (user *UserAdapter) UserSignup(req entities.User) (entities.User, error) {
	var res entities.User
	query := `INSERT INTO USERS (name,email,password)VALUES($1,$2,$3) RETURNING id,name,email`
	if err := user.DB.Raw(query, req.Name, req.Email, req.Password).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
