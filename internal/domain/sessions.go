package domain

import "time"

type Session struct {
	Id           string    `db:"id"`
	UserId       string    `db:"userid"`
	RefreshToken string    `db:"refreshtoken"`
	ExpiresIn    time.Time `db:"expiresin"`
	CreatedAt    time.Time `db:"createdat"`
}

type AdminSession struct {
	Id           string    `db:"id"`
	AdminId      string    `db:"adminid"`
	RefreshToken string    `db:"refreshtoken"`
	ExpiresIn    time.Time `db:"expiresin"`
	CreatedAt    time.Time `db:"createdat"`
}
