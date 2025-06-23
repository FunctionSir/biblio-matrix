/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 08:41:27
 * @LastEditTime: 2025-06-23 10:13:17
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/http.go
 */
package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

func apiOnlyHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a Biblio Matrix API-ONLY server."))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	passwd := r.PostFormValue("passwd")
	role := r.PostFormValue("role")
	if role == "" || (role != "admin" && role != "reader") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var token string
	var exp time.Time
	if role == "admin" {
		if username == "" || passwd == "" || !AuthAdmin(username, passwd) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token, exp = NewToken(username, true)
	} else {
		if username == "" || passwd == "" || !AuthReader(username, passwd) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token, exp = NewToken(username, false)
	}
	http.SetCookie(w, &http.Cookie{Name: "token", Value: token, SameSite: http.SameSiteDefaultMode, Expires: exp})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"status": "success", "admin": ChkTokensIsAdmin(token)})
}

func deauthHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	DelToken(tokenCookie.Value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func borrowHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	book := r.PostFormValue("book")
	duration := r.PostFormValue("duration")
	if book == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if duration == "" {
		duration = "30"
	}
	durationInt, err := strconv.Atoi(duration)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	result := Borrow(ctx, GetTokenUsername(tokenCookie.Value), book, time.Now().UTC(), time.Now().UTC().Add(time.Duration(durationInt*24)*time.Hour))
	if ctx.Err() != nil {
		http.Error(w, "操作超时. 请联系管理员.", http.StatusConflict)
		return
	}
	if result != "" {
		http.Error(w, result, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("🎉 恭喜! 借书成功!"))
}

func returnHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	book := r.PostFormValue("book")
	if book == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := Return(ctx, GetTokenUsername(tokenCookie.Value), book)
	if ctx.Err() != nil {
		http.Error(w, "操作超时. 请联系管理员.", http.StatusConflict)
		return
	}
	if result != "" {
		http.Error(w, result, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("🎉 恭喜! 还书成功!"))
}

func listBooksHandler(w http.ResponseWriter, r *http.Request) {
	books := ListBooks()
	if books == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func listRecordsHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	if !ChkTokensIsAdmin(tokenCookie.Value) && GetTokenUsername(tokenCookie.Value) != username {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if username == "" {
		username = "*"
	}
	records := ListRecords(username)
	if records == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	book := r.PostFormValue("book")
	if book == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	exists, err := IsBookExists(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	count := r.PostFormValue("count")
	if count == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	cnt, err := strconv.Atoi(count)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if exists {
		AddCnt(book, cnt)
		return
	}
	name := r.PostFormValue("name")
	author := r.PostFormValue("author")
	priceStr := r.PostFormValue("price")
	priceFloat64, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	priceInt := int(math.Round(priceFloat64 * 100))
	countStr := r.PostFormValue("count")
	countInt, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if name == "" || author == "" || priceInt < 0 || countInt <= 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = AddBook(Book{Id: book, Name: name, Author: author, Price: priceInt, Count: countInt})
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func newReaderHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	passwd := r.PostFormValue("passwd")
	name := r.PostFormValue("name")
	if username == "" || passwd == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = AddReader(username, passwd, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func newAdminHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	passwd := r.PostFormValue("passwd")
	name := r.PostFormValue("name")
	if username == "" || passwd == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = AddAdmin(username, passwd, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func delUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := DelUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func delBookHandler(w http.ResponseWriter, r *http.Request) {
	book := r.PostFormValue("book")
	if book == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := DelBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func listOverdueReadersHandler(w http.ResponseWriter, r *http.Request) {
	readers, err := ListOverdueReaders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(readers)
}

func readerInfoHandler(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Valid() != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !ChkTokensIsAdmin(tokenCookie.Value) && GetTokenUsername(tokenCookie.Value) != username {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	reader, err := GetReaderInfo(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reader)
}

func adminInfoHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	if username == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	admin, err := GetAdminInfo(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(admin)
}

func bookInfoHandler(w http.ResponseWriter, r *http.Request) {
	bookId := r.PostFormValue("book")
	if bookId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	book, err := GetBookInfo(bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func clearTokens(w http.ResponseWriter, r *http.Request) {

}

func serveHttp(addr string) {
	if FrontendDir == "" {
		log.Println("No frontend specified, started as an API-ONLY server.")
		http.HandleFunc("/", apiOnlyHomeHandler)
	} else {
		http.Handle("/", http.FileServer(http.Dir(FrontendDir)))
	}
	http.HandleFunc("/auth", Chain(authHandler, Logging))
	http.HandleFunc("/deauth", Chain(deauthHandler, ReaderLvlAuth, Logging))
	http.HandleFunc("/borrow", Chain(borrowHandler, ReaderLvlAuth, Logging))
	http.HandleFunc("/return", Chain(returnHandler, ReaderLvlAuth, Logging))
	http.HandleFunc("/list/books", Chain(listBooksHandler, Logging))
	http.HandleFunc("/list/records", Chain(listRecordsHandler, ReaderLvlAuth, Logging))
	http.HandleFunc("/list/overdue", Chain(listOverdueReadersHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/add", Chain(addHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/new/reader", Chain(newReaderHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/new/admin", Chain(newAdminHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/clear/tokens", Chain(clearTokens, AdminLvlAuth, Logging))
	http.HandleFunc("/del/user", Chain(delUserHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/del/book", Chain(delBookHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/readerinfo", Chain(readerInfoHandler, ReaderLvlAuth, Logging))
	http.HandleFunc("/admininfo", Chain(adminInfoHandler, AdminLvlAuth, Logging))
	http.HandleFunc("/bookinfo", Chain(bookInfoHandler, Logging))
	var err error
	if TlsCert == "" || TlsKey == "" {
		log.Printf("Listening HTTP on %s...\n", addr)
		err = http.ListenAndServe(addr, nil)
	} else {
		log.Printf("Using cert file: %s\n", TlsCert)
		log.Printf("Using key file: %s\n", TlsKey)
		log.Printf("Listening HTTPS on %s...\n", addr)
		err = http.ListenAndServeTLS(addr, TlsCert, TlsKey, nil)
	}
	panic(err)
}
