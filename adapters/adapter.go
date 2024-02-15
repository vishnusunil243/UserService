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
func (user *UserAdapter) GetUser(id int) (entities.User, error) {
	var res entities.User
	query := `SELECT * FROM USERS WHERE id=?`
	if err := user.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
func (user *UserAdapter) GetAdmin(id int) (entities.Admin, error) {
	var res entities.Admin
	query := `SELECT * FROM admins WHERE id=?`
	if err := user.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return entities.Admin{}, err
	}
	return res, nil
}
func (admin *UserAdapter) AdminLogin(email string) (entities.Admin, error) {
	var res entities.Admin
	query := `SELECT * FROM admins WHERE email=?`
	if err := admin.DB.Raw(query, email).Scan(&res).Error; err != nil {
		return entities.Admin{}, err
	}
	return res, nil
}
func (user *UserAdapter) UserLogin(email string) (entities.User, error) {
	var res entities.User
	query := `SELECT * FROM users WHERE email=?`
	if err := user.DB.Raw(query, email).Scan(&res).Error; err != nil {
		return entities.User{}, err
	}
	return res, nil
}
func (sup *UserAdapter) SuperAdminLogin(email string) (entities.SuperAdmin, error) {
	var res entities.SuperAdmin
	query := `SELECT * FROM super_admins WHERE email=?`
	if err := sup.DB.Raw(query, email).Scan(&res).Error; err != nil {
		return entities.SuperAdmin{}, err
	}
	return res, nil
}
func (sup *UserAdapter) AddAdmin(req entities.Admin) (entities.Admin, error) {
	var res entities.Admin
	var lastId int
	sup.DB.Raw("select MAX(id) FROM admins").Scan(&lastId)
	query := `INSERT INTO admins (name,email,password,id) VALUES ($1,$2,$3,$4) RETURNING id,name,email`
	if err := sup.DB.Raw(query, req.Name, req.Email, req.Password, lastId+1).Scan(&res).Error; err != nil {
		return entities.Admin{}, err
	}
	return res, nil
}
func (sup *UserAdapter) GetAllAdmins() ([]entities.Admin, error) {
	var res []entities.Admin
	query := `SELECT * FROM admins`
	if err := sup.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.Admin{}, err
	}
	return res, nil
}
func (admin *UserAdapter) GetAllUsers() ([]entities.User, error) {
	var res []entities.User
	query := `SELECT * FROM users`
	if err := admin.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.User{}, err
	}
	return res, nil
}
func (address *UserAdapter) AddAddress(req entities.Address) error {
	insertQuery := `INSERT INTO addresses (user_id,city,district,state,road) VALUES ($1,$2,$3,$4,$5)`
	if err := address.DB.Exec(insertQuery, req.UserId, req.City, req.District, req.State, req.Road).Error; err != nil {
		return err
	}
	return nil
}
func (address *UserAdapter) RemoveAddress(userId int) error {
	deleteQuery := `DELETE FROM addresses WHERE user_id=?`
	if err := address.DB.Exec(deleteQuery, userId).Error; err != nil {
		return err
	}
	return nil
}
func (address *UserAdapter) GetAddress(userId int) (entities.Address, error) {
	sQuery := `SELECT * FROM addresses WHERE user_id=?`
	var res entities.Address
	if err := address.DB.Raw(sQuery, userId).Scan(&res).Error; err != nil {
		return entities.Address{}, err
	}
	return res, nil
}
