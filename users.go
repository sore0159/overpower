package main

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

const PWFILE = DATADIR + "passwords.txt"

var USERSPWS map[string]string

func LoadUserData() (map[string]string, error) {
	pwdfile, err := os.Open(PWFILE)
	if err != nil {
		Log(err)
		return nil, err
	}
	defer pwdfile.Close()
	m := map[string]string{}
	pwdfileText := bufio.NewScanner(pwdfile)
	for pwdfileText.Scan() {
		text := pwdfileText.Text()
		line := strings.Split(text, "=")
		if len(line) != 2 || strings.TrimSpace(line[0]) == "" || strings.TrimSpace(line[1]) == "" {
			Log("Bad PWFILE line:", text)
			continue
		}
		m[strings.TrimSpace(line[0])] = strings.TrimSpace(line[1])
	}
	if err = pwdfileText.Err(); err != nil {
		Log(err)
		return nil, err
	}
	return m, nil
}

func ValidLogin(userName, password string) bool {
	pw, ok := USERSPWS[userName]
	return ok && pw == password
}

func CreateUser(userName, password string) error {
	f, err := os.OpenFile(PWFILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		Log(err)
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(userName + " = " + password + "\n"); err != nil {
		Log(err)
		return err
	}
	USERSPWS[userName] = password
	return nil
}

func ValidPassword(password string) bool {
	if password == "" {
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
	for test, _ := range USERSPWS {
		if name == strings.ToLower(test) {
			return true
		}
	}
	reserved := []string{"turn", "admin", "mule", "login", "logout"}
	for _, test := range reserved {
		if name == test {
			return true
		}
	}
	return false
}

func ValidUserName(username string) bool {
	if username == "" {
		return false
	}
	for _, rn := range username {
		if !unicode.In(rn, unicode.L, unicode.N) {
			return false
		}
	}
	return true
}
