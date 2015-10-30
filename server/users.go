package main

import (
	"database/sql"
	"strings"
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
	return ValidText(password)
}

func UserNameAvail(username string) bool {
	if !ValidUserName(username) {
		return false
	}
	name := strings.ToLower(username)
	reserved := []string{"home", "planet", "yours", "static", "turn", "admin", "themule", "mule", "login", "logout"}
	for _, test := range reserved {
		if name == test {
			return false
		}
	}
	query := "SELECT name FROM userinfo WHERE lower(name) = $1"
	var found string
	err := USERDB.QueryRow(query, username).Scan(&found)
	if err == nil {
		return false
	}
	if err == sql.ErrNoRows {
		return true
	}
	Log(err)
	return false
}

func UserExists(userN string) bool {
	if !ValidUserName(userN) {
		return false
	}
	query := "SELECT name FROM userinfo WHERE lower(name) = $1"
	name := strings.ToLower(userN)
	var found string
	err := USERDB.QueryRow(query, name).Scan(&found)
	if err == nil && found != "" {
		return true
	}
	return false
}

func ValidUserName(username string) bool {
	return ValidText(username)
}
