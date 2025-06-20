/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 08:41:27
 * @LastEditTime: 2025-06-20 11:26:31
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/http.go
 */
package main

import (
	"log"
	"net/http"
)

func apiOnlyHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a Biblio Matrix API-ONLY server."))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	passwd := r.PostFormValue("passwd")
	if username == "" || passwd == "" || !Auth(username, passwd) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	token, exp := NewToken(ChkIsAdmin(username))
	http.SetCookie(w, &http.Cookie{Name: "token", Value: token, SameSite: http.SameSiteDefaultMode, Expires: exp})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusAccepted)))
}

func deauthHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("TOKEN")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	DelToken(tokenCookie.Value)
}

func serveHttp(addr string) {
	if FrontendDir == "" {
		log.Println("No frontend specified, started as an API-ONLY server.")
		http.HandleFunc("/", Chain(apiOnlyHomeHandler, Logging))
	} else {
		http.Handle("/", http.FileServer(http.Dir(FrontendDir)))
	}
	http.HandleFunc("/auth", Chain(authHandler, Logging))
	http.HandleFunc("/deauth", Chain(deauthHandler, SimpleAuth, Logging))
	err := http.ListenAndServe(addr, nil)
	panic(err)
}
