package domain

import "time"

type User struct {
	Id            string    `json:"_" db:"id"`
	Email         string    `json:"email" db:"email"`
	Password_hash string    `json:"password_hash" db:"password_hash"`
	First_name    string    `json:"first_name" db:"first_name"`
	Last_name     string    `json:"last_name" db:"last_name"`
	Registered_at time.Time `json:"registered_at" db:"registered_at"`
}
