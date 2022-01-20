package models

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	PhoneNo int    `json:"phone_no"`
}
