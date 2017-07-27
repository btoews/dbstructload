package main

import (
	"database/sql"
	"fmt"

	"github.com/mastahyeti/dbstructload"
)

var (
	db *sql.DB
)

// User represents a row from the `users` table.
type User struct {
	ID    uint64 `queryField:"User_id"`
	Login string `queryField:"User_login"`
	Email *UserEmail
}

// UserEmail represents a row from the `user_emails` table.
type UserEmail struct {
	ID     uint64 `queryField:"UserEmail_id"`
	UserID uint64 `queryField:"UserEmail_user_id"`
	Email  string `queryField:"UserEmail_email"`
}

func findUsersByLoginWithEmail(logins ...string) ([]*User, error) {
	const query = `
		SELECT
			users.id            AS User_id,
			users.login         AS User_login,
			user_emails.id      AS UserEmail_id,
			user_emails.user_id AS UserEmail_user_id,
			user_emails.email   AS UserEmail_email,
		FROM users
		JOIN user_emails
			ON user_emails.user_id = users.id
		WHERE login=?;
	`

	// Allocate space, assuming a user with every login exists.
	users := make([]*User, 0, len(logins))

	// Convert string slice to interface{} slice.
	queryArgs := make([]interface{}, len(logins))
	for i := range logins {
		queryArgs[i] = logins[i]
	}

	rows, err := dbstructload.Query(db, query, queryArgs)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := User{}
		user.Email = &UserEmail{}

		if err = rows.Load(&user, user.Email); err != nil {
			return users, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func main() {
	var err error

	if db, err = sql.Open("mysql", "someaddress"); err != nil {
		panic(err)
	}

	users, err := findUsersByLoginWithEmail("userone", "usertwo", "userthree")
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		fmt.Printf("%s: %s\n", user.Login, user.Email.Email)
	}
}
