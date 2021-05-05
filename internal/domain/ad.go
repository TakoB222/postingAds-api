package domain

import "github.com/lib/pq"

type (
	Ad struct {
		Id          int            `db:"id"`
		UserId      string         `db:"userid"`
		Title       string         `db:"title"`
		Category    string         `db:"category_id"`
		Description string         `db:"description"`
		Price       int            `db:"price"`
		Contacts    string         `db:"contacts_id"`
		Published   bool           `db:"published"`
		ImagesURL   pq.StringArray `db:"images_url"`
	}

	Contacts struct {
		Name         string `db:"name"`
		Phone_number string `db:"phone_number"`
		Email        string `db:"email"`
		Location     string `db:"location"`
	}

	Categories struct {
		Id             int    `json:"id"`
		Category       string `json:"category"`
		ParentCategory int    `json:"parent_category"`
	}
)
