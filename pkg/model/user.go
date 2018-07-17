package model

// User struct
type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Client   Client `gorm:"save_associations:false"`
	ClientID uint   `gorm:"unique"`
}
