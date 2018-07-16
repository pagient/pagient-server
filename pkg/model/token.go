package model

// Token struct
type Token struct {
	Token string `json:"token"`
	User  string `json:"-"`
}
