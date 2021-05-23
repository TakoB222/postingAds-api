package repository

import (
	"errors"
	"fmt"
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

type AdRepository struct {
	db *sqlx.DB
}

func NewAdRepository(db *sqlx.DB) *AdRepository {
	return &AdRepository{db: db}
}

func (r *AdRepository) GetAllAdsByUserId(userId string) ([]domain.Ad, error) {
	var ads []domain.Ad

	tx, err := r.db.Begin()
	if err != nil {
		return []domain.Ad{}, err
	}
	query := fmt.Sprintf("select * from %s where userid=$1", database.AdsTable)
	if err := r.db.Select(&ads, query, userId); err != nil {
		return nil, err
	}

	for i := 0; i < len(ads); i++ {
		categoryId, err := strconv.Atoi(ads[i].Category)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return []domain.Ad{}, err
			}
			return []domain.Ad{}, err
		}

		var categorySeq []string
		query = fmt.Sprintf(`with recursive r as (select id, parent_category, category from categories where id=$1 
									union 
									select categories.id, categories.parent_category, categories.category from %s join r on categories.id = r.parent_category) 
									select category from r;`, database.CategoriesTable)
		if err := r.db.Select(&categorySeq, query, categoryId); err != nil {
			err := tx.Rollback()
			if err != nil {
				return []domain.Ad{}, err
			}
			return []domain.Ad{}, err
		}

		ads[i].Category = func(categories []string) string{
			for i := 0; i < len(categories)/2; i++ {
				tmp := categories[i]
				categories[i] = categories[len(categories)-1-i]
				categories[len(categories)-1-i] = tmp
			}
			return strings.Join(categories, "/")
		}(categorySeq)
	}
	if err = tx.Commit(); err != nil {
		return []domain.Ad{}, err
	}

	return ads, nil
}

func (r *AdRepository) CreateAd(userId string, input Ads) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var category int
	query := fmt.Sprintf("select id from %s where category=$1", database.CategoriesTable)
	err = r.db.Get(&category, query, input.Category)
	if err != nil {
		return 0, err
	}

	var contactId int
	query = fmt.Sprintf("insert into %s (name, phone_number, email, location) values ($1, $2, $3, $4) returning id", database.ContactsInfoTable)
	row := tx.QueryRow(query, input.Contacts.Name, input.Contacts.Phone_number, input.Contacts.Email, input.Contacts.Location)
	if err := row.Scan(&contactId); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	var adId int
	query = fmt.Sprintf("insert into %s (userid, title, category_id, description, price, contacts_id, published, images_url) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id", database.AdsTable)
	row = tx.QueryRow(query, userId, input.Title, category, input.Description, input.Price, contactId, input.Published, pq.Array(input.ImagesURL))
	if err := row.Scan(&adId); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	return adId, tx.Commit()
}

func (r *AdRepository) GetAdById(userId string, adId string) (domain.Ad, error) {
	var ad domain.Ad

	tx, err := r.db.Begin()
	if err != nil {
		return domain.Ad{}, err
	}

	query := fmt.Sprintf("select * from %s where userid=$1 and id=$2", database.AdsTable)
	if err := r.db.Get(&ad, query, userId, adId); err != nil {
		return domain.Ad{}, err
	}

	categoryId, err := strconv.Atoi(ad.Category)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return domain.Ad{}, err
		}
		return domain.Ad{}, err
	}

	var categorySeq []string
	query = fmt.Sprintf(`with recursive r as (select id, parent_category, category from categories where id=$1 
								union 
								select categories.id, categories.parent_category, categories.category from %s join r on categories.id = r.parent_category) 
								select category from r;`, database.CategoriesTable)
	if err := r.db.Select(&categorySeq, query, categoryId); err != nil {
		err := tx.Rollback()
		if err != nil {
			return domain.Ad{}, err
		}
		return domain.Ad{}, err
	}

	ad.Category = func(categories []string) string{
		for i := 0; i < len(categories)/2; i++ {
			tmp := categories[i]
			categories[i] = categories[len(categories)-1-i]
			categories[len(categories)-1-i] = tmp
		}
		return strings.Join(categories, "/")
	}(categorySeq)

	if err = tx.Commit(); err != nil {
		return domain.Ad{}, err
	}

	return ad, nil
}

func (r *AdRepository) UpdateAd(userId string, adId string, ad Ads) error {
	var categoryId int
	query := fmt.Sprintf("select id from %s where category=$1", database.CategoriesTable)
	if err := r.db.Get(&categoryId, query, ad.Category); err != nil {
		return errors.New("here")
	}

	var contactsId int
	query = fmt.Sprintf("select contacts_id from %s where id=$1 and userid=$2", database.AdsTable)
	if err := r.db.Get(&contactsId, query, adId, userId); err != nil {
		return err
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if ad.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, ad.Title)
		argId++
	}

	if ad.Category != "" {
		setValues = append(setValues, fmt.Sprintf("category_id=$%d", argId))
		args = append(args, categoryId)
		argId++
	}

	if ad.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, ad.Description)
		argId++
	}

	if ad.Price >= 0 {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, ad.Price)
		argId++
	}

	if len(ad.ImagesURL) > 0 {
		setValues = append(setValues, fmt.Sprintf("images_url=$%d", argId))
		args = append(args, pq.Array(ad.ImagesURL))
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("published=$%d", argId))
	args = append(args, ad.Published)
	argId++

	tx, err := r.db.Begin()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	query = fmt.Sprintf("update %s set name=$1, phone_number=$2, email=$3, location=$4 where id=$5", database.ContactsInfoTable)
	if _, err := tx.Exec(query, ad.Contacts.Name, ad.Contacts.Phone_number, ad.Contacts.Email, ad.Contacts.Location, contactsId); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	setValues = append(setValues, fmt.Sprintf("contacts_id=$%d", argId))
	args = append(args, contactsId)
	argId++

	setQuery := strings.Join(setValues, ", ")

	query = fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d and userid = $%d",
		database.AdsTable, setQuery, argId, argId+1)
	args = append(args, adId, userId)

	if _, err = r.db.Exec(query, args...); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r *AdRepository) DeleteAd(userId string, adId string) error {
	query := fmt.Sprintf("delete from %s where id=$1 and userid=$2", database.AdsTable)
	if _, err := r.db.Exec(query, adId, userId); err != nil {
		return err
	}

	return nil
}

func (r *AdRepository) SearchAdByRequest(search_request string)([]FtsResponse, error) {
	var res []FtsResponse

	query := fmt.Sprintf("select id, ts_headline(title, q) as title from %s, plainto_tsquery('russian', $1) as q where make_tsvector(title, description) @@ q order by ts_rank(make_tsvector(title, description), q) desc", database.AdsTable)
	if err := r.db.Select(&res, query, search_request); err != nil {
		return nil, err
	}

	return res, nil
}
