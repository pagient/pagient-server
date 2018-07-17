package model

// Pager struct
type Pager struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null;unique"`
}
