package main

import (
	//"mule/overpower"
	"net/http"
)

func SetupMux() {
	http.HandleFunc("/", pageMainIndex)
	http.HandleFunc("/favicon.ico", imgFavIcon)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))

	http.HandleFunc("/auth/", pageAuthIndex)
	http.HandleFunc("/auth/logout", pageAuthLogout)
	http.HandleFunc("/auth/login", pageAuthLogin)
	http.HandleFunc("/auth/create", pageAuthCreate)

	http.HandleFunc("/overpower/", pageOPIndex)
	http.HandleFunc("/overpower/list/", pageOPList)
	http.HandleFunc("/overpower/view/", pageOPView)
	http.HandleFunc("/overpower/home", pageOPHome)
	//http.HandleFunc("/overpower/quit/", pageOPQuit)
	//http.HandleFunc("/overpower/play/", pageOPPlay)

	//http.HandleFunc("/overpower/json/", apiJson)

}

func imgFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "STATIC/img/favicon.ico")
}
