package models

type User struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Links    []*Links `json:"links"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
