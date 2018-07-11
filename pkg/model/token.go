package model

type Token struct {
	Token string `json:"token"`
	User  string `json:"-"`
}
