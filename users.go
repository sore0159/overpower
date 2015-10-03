package main

import (
	"database/sql"
	"strings"
	"unicode"
)

func ValidLogin(userName, password string) bool {
	var pass string
	err := USERDB.QueryRow("SELECT password FROM userinfo WHERE name = $1", userName).Scan(&pass)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		Log(err)
		return false
	}
	if pass == "" {
		return false
	}
	return pass == password
}

func CreateUser(userName, password string) error {
	query := "INSERT INTO userinfo (name, password) VALUES($1, $2)"
	res, err := USERDB.Exec(query, userName, password)
	if err != nil {
		return Log(err)
	}
	if aff, err := res.RowsAffected(); err != nil {
		return Log(err)
	} else if aff == 0 {
		return Log("Database update failed")
	}
	return nil
}

func ValidPassword(password string) bool {
	if password == "" || len(password) > 15 {
		return false
	}
	for _, rn := range password {
		if !unicode.In(rn, unicode.L, unicode.N) {
			return false
		}
	}
	return true
}

func UserNameInUse(username string) bool {
	name := strings.ToLower(username)
	reserved := []string{"planet", "yours", "static", "turn", "admin", "mule", "login", "logout"}
	for _, test := range reserved {
		if name == test {
			return true
		}
	}
	query := "SELECT name FROM userinfo WHERE lower(name) = $1"
	var found string
	err := USERDB.QueryRow(query, username).Scan(&found)
	if err == nil {
		return true
	}
	if err == sql.ErrNoRows {
		return false
	}
	Log(err)
	return true
}

func ValidUserName(username string) bool {
	if username == "" || len(username) > 15 {
		return false
	}
	for _, rn := range username {
		if !unicode.In(rn, unicode.L, unicode.N) {
			return false
		}
	}
	return true
}
