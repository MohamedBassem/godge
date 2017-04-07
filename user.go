package godge

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func (u *user) save(db *sqlx.DB) error {
	_, err := db.NamedExec("INSERT INTO users (username, password) VALUES (:username, :password)", u)
	return err
}

func (u *user) isCorrectPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

var userQ userQuery = userQuery{}

type userQuery struct{}

func (*userQuery) find(db *sqlx.DB, username string) (*user, error) {
	u := &user{}
	if err := db.Get(u, "SELECT * FROM users WHERE username=?", username); err != nil {
		return nil, err
	}
	return u, nil
}

func (*userQuery) findAll(db *sqlx.DB) ([]user, error) {
	u := []user{}
	if err := db.Get(&u, "SELECT * FROM users"); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *userQuery) usernames(db *sqlx.DB) ([]string, error) {
	us, err := u.findAll(db)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, k := range us {
		ret = append(ret, k.Username)
	}
	return ret, nil
}
