package model

// User struct
type User struct {
	Username string  `json:"username"`
	Password string  `json:"password,omitempty"`
	Client   *Client `json:"-"`
}
