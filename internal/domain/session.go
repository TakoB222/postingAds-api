package domain

import "time"

type Session struct {
	Id           string    `db:"id"`
	UserId       string    `db:"userid"`
	RefreshToken string    `db:"refreshtoken"`
	UA           string    `db:"ua"`
	Ip           string    `db:"ip"`
	ExpiresIn    time.Time `db:"expiresin"`
	CreatedAt    time.Time `db:"createdat"`
}
