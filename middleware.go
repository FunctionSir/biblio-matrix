package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/FunctionSir/goset"
	"github.com/google/uuid"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { log.Println(r.RemoteAddr, r.URL.Path, time.Since(start)) }()
		next(w, r)
	}
	return handler
}

var TokensSet goset.Set[string]
var TokensExp map[string]time.Time
var TokensUser map[string]string
var TokensIsAdmin map[string]bool
var TokenLock sync.Mutex

func ChkToken(token string) bool {
	TokenLock.Lock()
	defer TokenLock.Unlock()
	if !TokensSet.Has(token) {
		return false
	}
	// Do NOT use DelToken here!
	// Or a DEAD LOCK will occur!
	if TokensExp[token].UnixNano() <= time.Now().UnixNano() {
		TokensSet.Erase(token)
		delete(TokensExp, token)
		delete(TokensUser, token)
		delete(TokensIsAdmin, token)
		return false
	}
	return true
}

func ChkTokensIsAdmin(token string) bool {
	TokenLock.Lock()
	defer TokenLock.Unlock()
	if !TokensSet.Has(token) {
		return false
	}
	return TokensIsAdmin[token]
}

func GetTokenUsername(token string) string {
	TokenLock.Lock()
	defer TokenLock.Unlock()
	if !TokensSet.Has(token) {
		return ""
	}
	return TokensUser[token]
}

func NewToken(username string, isAdmin bool) (string, time.Time) {
	TokenLock.Lock()
	defer TokenLock.Unlock()
	token := uuid.NewString()
	for TokensSet.Has(token) {
		token = uuid.NewString()
	}
	exp := time.Now().Add(24 * time.Hour)
	TokensSet.Insert(token)
	TokensExp[token] = exp
	TokensUser[token] = username
	TokensIsAdmin[token] = isAdmin
	return token, exp
}

func DelToken(token string) {
	TokenLock.Lock()
	defer TokenLock.Unlock()
	if !TokensSet.Has(token) {
		return
	}
	TokensSet.Erase(token)
	delete(TokensExp, token)
	delete(TokensUser, token)
	delete(TokensIsAdmin, token)
}

func ReaderLvlAuth(next http.HandlerFunc) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Valid() != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if !ChkToken(tokenCookie.Value) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
	return handler
}

// Do not need to use with SimpleAuth, just use it ONLY when needed.
func AdminLvlAuth(next http.HandlerFunc) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil || tokenCookie.Valid() != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if !ChkToken(tokenCookie.Value) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if !ChkTokensIsAdmin(tokenCookie.Value) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next(w, r)
	}
	return handler
}

func Chain(handlerfunc http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handlerfunc = middleware(handlerfunc)
	}
	return handlerfunc
}
