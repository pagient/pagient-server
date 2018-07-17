package model

// Client struct
type Client struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null;unique"`
}
