package model

type Person struct {
	Name  string `json:"name" db:"name"`
	IIN   string `json:"iin" db:"iin"`
	Phone string `json:"phone" db:"phone"`
}

type IINResponse struct {
	Correct     bool   `json:"correct"`
	Sex         string `json:"sex,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

type PersonResponse struct {
	Success bool   `json:"success"`
	Errors  string `json:"errors,omitempty"`
}
