package entities

type User struct {
	Id       uint
	Name     string
	Password string
	Email    string
}
type Admin struct {
	Id       uint
	Name     string
	Email    string
	Password string
}
type SuperAdmin struct {
	Id       uint
	Name     string
	Email    string
	Password string
}
type Address struct {
	Id       uint
	UserId   uint
	User     `gorm:"foreignKey:UserId"`
	City     string
	District string
	State    string
	Road     string
}
