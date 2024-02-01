package auth

import (
	"database/sql"
	"forum/internal/model"
)

type AuthRepo struct {
	DB *sql.DB
}

type Authorization interface {
	CreateUser(user model.User) error
	GetUserByUsername(username string) (model.User, error)
	CreateSession(session model.Session) error
	DeleteSessionByUserId(UserId int) error
	DeleteSessionByToken(token string) error
	GetUserBySession(token string) (model.User, error)
}

func NewAuthorizationRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		DB: db,
	}
}

// CRITICAL ERROR: AUTH SERVICE DOES NOT RECORD USERID + IT DELETES ALL RECORD ABOUT OTHER SESSION

func (a *AuthRepo) CreateUser(u model.User) error {
	query := `Insert into users (Username, Email, Password, AuthMode) Values ($1, $2, $3, $4)`

	_, err := a.DB.Exec(query, u.Username, u.Email, u.Password, u.AuthMethod)
	if err != nil {
		return err
	}

	return err
}

func (a *AuthRepo) CreateSession(session model.Session) error {
	query := "INSERT INTO Session (Token, ExpDate, UserId) VALUES ($1, $2, $3)"

	_, err := a.DB.Exec(query, session.Token, session.ExpTime, session.UserId)

	return err
}

func (a *AuthRepo) DeleteSessionByUserId(UserId int) error {
	query := "DELETE FROM Session WHERE UserId = $1"
	res, err := a.DB.Exec(query, UserId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *AuthRepo) DeleteSessionByToken(token string) error {
	query := "DELETE FROM Session WHERE Token = $1"
	_, err := a.DB.Exec(query, token)

	return err
}

func (a *AuthRepo) GetUserByUsername(username string) (model.User, error) {
	query := "SELECT * FROM Users WHERE Username = $1"

	Result := model.User{}

	err := a.DB.QueryRow(query, username).Scan(&Result.Id, &Result.Email, &Result.Username, &Result.Password, &Result.AuthMethod)

	return Result, err
}

func (a *AuthRepo) GetUserBySession(token string) (model.User, error) {
	var user model.User

	query := "SELECT UserId, Username FROM Session WHERE Token = $1 AND ExpDate > CURRENT_TIMESTAMP"

	query2 := "SELECT Users.Id, Users.Email, Users.Username FROM Users INNER JOIN Session ON Users.Id = Session.UserId WHERE Session.Token = $1 AND Session.ExpDate > CURRENT_TIMESTAMP"

	_ = query

	err := a.DB.QueryRow(query2, token).Scan(&user.Id, &user.Email, &user.Username)

	return user, err
}
