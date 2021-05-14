package repository

import (
	"fmt"
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"strings"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAdminId(email, password_hash string) (string, error) {
	var id int

	query := fmt.Sprintf("select id from %s where login=$1 and password_hash=$2", database.AdminsTable)
	if err := r.db.Get(&id, query, email, password_hash); err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

func (r *AdminRepository) SetAdminSession(session domain.AdminSession) error {
	query := fmt.Sprintf("insert into %s (adminid, refreshtoken, expiresin, createdat) values ($1, $2, $3, $4)", database.AdminRefreshSessionTable)
	if _, err := r.db.Exec(query, session.AdminId, session.RefreshToken, session.ExpiresIn, session.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) GetAdminSessionByRefreshToken(refrehsToken string) (domain.AdminSession, error) {
	var session domain.AdminSession
	query := fmt.Sprintf("select * from %s where refreshtoken=$1", database.AdminRefreshSessionTable)
	if err := r.db.Get(&session, query, refrehsToken); err != nil {
		return domain.AdminSession{}, err
	}

	return session, nil
}

func (r *AdminRepository) DeleteAdminSessionByAdminId(adminId string) error {
	query := fmt.Sprintf("delete from %s where adminid=$1", database.AdminRefreshSessionTable)
	if _, err := r.db.Exec(query, adminId); err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) GetAllAdsByAdmin() ([]domain.Ad, error) {
	var ads []domain.Ad

	query := fmt.Sprintf("select * from %s", database.AdsTable)
	if err := r.db.Select(&ads, query); err != nil {
		return nil, err
	}

	return ads, nil
}

func (r *AdminRepository) GetAd(adId string) (domain.Ad, error) {
	var ad domain.Ad

	query := fmt.Sprintf("select * from %s where id=$1", database.AdsTable)
	if err := r.db.Get(&ad, query, adId); err != nil {
		return domain.Ad{}, err
	}

	return ad, nil
}

func (r *AdminRepository) AdminDeleteAd(adId string) error {
	query := fmt.Sprintf("delete from %s where id=$1", database.AdsTable)
	if _, err := r.db.Exec(query, adId); err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) AdminUpdateAd(adId string, ad Ads) error {
	var categoryId int
	query := fmt.Sprintf("select id from %s where category=$1", database.CategoriesTable)
	if err := r.db.Get(&categoryId, query, ad.Category); err != nil {
		return err
	}

	var userId string
	query = fmt.Sprintf("select userid from %s where id=$1", database.AdsTable)
	if err := r.db.Get(&userId, query, adId); err != nil {
		return err
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
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("update %s set name=$1, phone_number=$2, email=$3, location=$4 where id=$5", database.ContactsInfoTable)
	if _, err := tx.Exec(query, ad.Contacts.Name, ad.Contacts.Phone_number, ad.Contacts.Email, ad.Contacts.Location, contactsId); err != nil {
		tx.Rollback()
		return err
	}

	setValues = append(setValues, fmt.Sprintf("contacts_id=$%d", argId))
	args = append(args, contactsId)
	argId++

	setQuery := strings.Join(setValues, ", ")

	query = fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d and userid = $%d",
		database.AdsTable, setQuery, argId, argId+1)
	args = append(args, adId, userId)
	fmt.Println(query)

	if _, err = r.db.Exec(query, args...); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
