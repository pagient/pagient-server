package model

// Token struct
type Token struct {
	ID     uint   `gorm:"primary_key"`
	Raw    string `gorm:"not null;unique"`
	User   User   `gorm:"save_associations:false"`
	UserID uint
}
