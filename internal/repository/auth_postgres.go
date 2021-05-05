package repository

import (
	"fmt"
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/pkg/database"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (email, password_hash, first_name, last_name, registered_At) values ($1, $2, $3, $4, $5) returning id", database.UsersTable)

	row := r.db.QueryRow(query, user.Email, user.Password_hash, user.First_name, user.Last_name, user.Registered_at)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthRepository) GetUser(email, password_hash string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf("Select id from %s where email=$1 and password_hash=$2", database.UsersTable)

	err := r.db.Get(&user, query, email, password_hash)
	if err != nil {
		return domain.User{}, err
	}

	return user, err
}

func (r *AuthRepository) GetSessionByRefreshToken(refreshToken string) (domain.Session, error) {
	//TODO: if ua and ip wrong, what then...
	var session domain.Session

	query := fmt.Sprintf("select * from %s where refreshtoken=$1", database.RefreshSessionsTable)

	err := r.db.Get(&session, query, refreshToken)
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r *AuthRepository) DeleteSessionByUserId(userId string) error {
	query := fmt.Sprintf("delete from %s where userId=$1", database.RefreshSessionsTable)

	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) SetSession(session domain.Session) error {
	query := fmt.Sprintf("Insert into %s (userId, refreshToken, ua, ip, expiresIn, createdAt) values ($1,$2,$3,$4,$5,$6)", database.RefreshSessionsTable)
	_, err := r.db.Exec(query, session.UserId, session.RefreshToken, session.UA, session.Ip, session.ExpiresIn, session.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
